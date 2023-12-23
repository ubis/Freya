package packet

import (
	"bytes"
	"time"

	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/subpasswd"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"

	"golang.org/x/crypto/bcrypt"
)

// SubPasswordSet Packet
func SubPasswordSet(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	passwd := string(bytes.Trim(reader.ReadBytes(10), "\x00"))
	reader.ReadBytes(5)

	question := reader.ReadInt32()
	answer := string(bytes.Trim(reader.ReadBytes(16), "\x00"))

	pkt := network.NewWriter(CSCSubPasswordSet)

	sub := session.SubPassword

	if len(passwd) < 4 || question < 1 || question > 10 {
		pkt.WriteInt32(0x00) // failed
		pkt.WriteInt32(0x00)
		pkt.WriteInt32(0x01)
		pkt.WriteInt32(0x00)

		log.Info("a")
		network.DumpPacket(pkt)

		session.Send(pkt)
		return
	}

	if sub.Password == "" {
		// creating sub password for the first time
		// check answer
		if len(answer) < 4 {
			pkt.WriteInt32(0x00) // failed
			pkt.WriteInt32(0x00)
			pkt.WriteInt32(0x01)
			pkt.WriteInt32(0x00)

			log.Info("b")
			network.DumpPacket(pkt)

			session.Send(pkt)
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(answer), bcrypt.DefaultCost)
		sub.Answer = string(hash)
		sub.Question = byte(question)
		sub.Expires = time.Now()
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	sub.Password = string(hash)

	// update to db
	req := subpasswd.SetReq{Account: session.Account, Details: *sub}
	res := subpasswd.SetRes{}
	err := session.RPC.Call(rpc.SetSubPassword, &req, &res)

	if err == nil && res.Success {
		pkt.WriteInt32(0x01) // success
	} else {
		pkt.WriteInt32(0x00) // failed
	}

	pkt.WriteInt32(0x00)
	pkt.WriteInt32(0x01)
	pkt.WriteInt32(0x00)

	network.DumpPacket(pkt)

	session.Send(pkt)
}

// SubPasswordCheckRequest Packet
func SubPasswordCheckRequest(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	sub := session.SubPassword
	left := time.Until(sub.Expires)
	state := 1 // 1 - verification is needed; 0 - not needed

	pkt := network.NewWriter(CSCSubPasswordCheckRequest)

	if session.ServerConfig.IgnoreSubPassword || left.Seconds() > 0 {
		state = 0
	}

	pkt.WriteInt32(state)
	session.Send(pkt)
}

// SubPasswordCheck Packet
func SubPasswordCheck(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	password := string(bytes.Trim(reader.ReadBytes(10), "\x00"))
	reader.ReadBytes(5)
	hours := reader.ReadInt32()

	sub := session.SubPassword
	err := bcrypt.CompareHashAndPassword([]byte(sub.Password), []byte(password))

	pkt := network.NewWriter(CSCSubPasswordCheck)

	if hours < 0 || hours > 4 {
		pkt.WriteInt32(0x00) // failed
		pkt.WriteByte(sub.FailTimes)
		pkt.WriteInt32(0x00)
		pkt.WriteInt32(0x01)

		sub.FailTimes++
		session.Send(pkt)
		return
	}

	if err != nil {
		pkt.WriteInt32(0x00) // failed
		sub.FailTimes++
	} else {
		sub.Expires = time.Now().Add(time.Hour * time.Duration(hours))
		req := subpasswd.SetReq{Account: session.Account, Details: *sub}
		res := subpasswd.SetRes{}
		err := session.RPC.Call(rpc.SetSubPassword, &req, &res)

		if err != nil || !res.Success {
			pkt.WriteInt32(0x00) // failed
			sub.FailTimes++
		} else {
			pkt.WriteInt32(0x01) // success
			sub.FailTimes = 0
			sub.Verified = true
		}
	}

	pkt.WriteByte(sub.FailTimes)
	pkt.WriteInt32(0x00)
	pkt.WriteInt32(0x01)

	session.Send(pkt)
}

// SubPasswordFindRequest Packet
func SubPasswordFindRequest(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	sub := session.SubPassword

	pkt := network.NewWriter(CSCSubPasswordFindRequest)
	pkt.WriteInt32(sub.Question)
	pkt.WriteInt32(sub.Question)
	pkt.WriteInt32(0x01)

	session.Send(pkt)
}

// SubPasswordFind Packet
func SubPasswordFind(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	reader.ReadBytes(8)
	answer := string(bytes.Trim(reader.ReadBytes(16), "\x00"))

	sub := session.SubPassword
	err := bcrypt.CompareHashAndPassword([]byte(sub.Answer), []byte(answer))

	pkt := network.NewWriter(CSCSubPasswordFind)

	if err != nil {
		pkt.WriteInt32(0x00) // failed
		sub.FailTimes++
	} else {
		pkt.WriteInt32(0x01) // success
		sub.FailTimes = 0
		sub.Verified = true
	}

	pkt.WriteByte(sub.FailTimes)
	pkt.WriteInt32(0x01)

	session.Send(pkt)
}

// SubPasswordDeleteRequest Packet
func SubPasswordDeleteRequest(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	reader.ReadBytes(4)
	password := string(bytes.Trim(reader.ReadBytes(10), "\x00"))

	sub := session.SubPassword
	err := bcrypt.CompareHashAndPassword([]byte(sub.Password), []byte(password))

	pkt := network.NewWriter(CSCSubPasswordDeleteRequest)

	if err != nil {
		pkt.WriteInt32(0x00) // failed
		sub.FailTimes++
	} else {
		pkt.WriteInt32(0x01) // success
		sub.FailTimes = 0
		sub.Verified = true
	}

	pkt.WriteByte(sub.FailTimes)
	pkt.WriteInt32(0x01)

	session.Send(pkt)
}

// SubPasswordDelete Packet
func SubPasswordDelete(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	pkt := network.NewWriter(CSCSubPasswordDelete)

	sub := session.SubPassword

	if !sub.Verified {
		pkt.WriteInt32(0x00) // failed
		pkt.WriteInt32(0x01)

		session.Send(pkt)
		return
	}

	// update to db
	req := subpasswd.SetReq{Account: session.Account}
	res := subpasswd.SetRes{}
	err := session.RPC.Call(rpc.RemoveSubPassword, req, &res)

	if err == nil && res.Success {
		*sub = subpasswd.Details{}
		pkt.WriteInt32(0x01) // success
	} else {
		pkt.WriteInt32(0x00) // failed
	}

	pkt.WriteInt32(0x01)

	session.Send(pkt)
}

// SubPasswordChangeQARequest Packet
func SubPasswordChangeQARequest(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	reader.ReadBytes(4)
	password := string(bytes.Trim(reader.ReadBytes(10), "\x00"))

	sub := session.SubPassword
	err := bcrypt.CompareHashAndPassword([]byte(sub.Password), []byte(password))

	pkt := network.NewWriter(CSCSubPasswordChangeQARequest)

	if err != nil {
		pkt.WriteInt32(0x00) // failed
		sub.FailTimes++
	} else {
		pkt.WriteInt32(0x01) // success
		sub.FailTimes = 0
		sub.Verified = true
	}

	pkt.WriteByte(sub.FailTimes)
	pkt.WriteInt32(0x01)

	session.Send(pkt)
}

// SubPasswordChangeQA Packet
func SubPasswordChangeQA(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	reader.ReadBytes(4)
	question := reader.ReadInt32()
	answer := string(bytes.Trim(reader.ReadBytes(16), "\x00"))

	pkt := network.NewWriter(CSCSubPasswordChangeQA)

	sub := session.SubPassword

	if len(answer) < 4 || question < 1 || question > 10 || !sub.Verified {
		pkt.WriteInt32(0x00) // failed
		pkt.WriteInt32(0x01)

		session.Send(pkt)
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(answer), bcrypt.DefaultCost)
	sub.Answer = string(hash)
	sub.Question = byte(question)

	// update to db
	req := subpasswd.SetReq{Account: session.Account, Details: *sub}
	res := subpasswd.SetRes{}
	err := session.RPC.Call(rpc.SetSubPassword, req, &res)

	if err == nil && res.Success {
		pkt.WriteInt32(0x01) // success
	} else {
		pkt.WriteInt32(0x00) // failed
	}

	sub.Verified = false
	pkt.WriteInt32(0x01)

	session.Send(pkt)
}

// CharacterDeleteCheckSubPassword Packet
func CharacterDeleteCheckSubPassword(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	password := string(bytes.Trim(reader.ReadBytes(10), "\x00"))

	sub := session.SubPassword
	err := bcrypt.CompareHashAndPassword([]byte(sub.Password), []byte(password))

	pkt := network.NewWriter(CSCCharacterDeleteCheckSubPassword)

	if err != nil {
		pkt.WriteInt32(0x00) // failed
		sub.FailTimes++
	} else {
		pkt.WriteInt32(0x01) // success
		sub.FailTimes = 0
		sub.Verified = true
	}

	pkt.WriteByte(sub.FailTimes)

	session.Send(pkt)
}
