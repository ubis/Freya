package packet

import (
	"bytes"
	"time"

	"github.com/ubis/Freya/cmd/loginserver/server"
	"github.com/ubis/Freya/share/event"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/account"
	"github.com/ubis/Freya/share/models/message"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

// PublicKey Packet
func PublicKey(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified) {
		return
	}

	rsa := session.ServerInstance.RSA
	key := rsa.PublicKey

	pkt := network.NewWriter(CSCPublicKey)
	pkt.WriteBool(true)
	pkt.WriteUint16(len(key))
	pkt.WriteBytes(key[:])

	session.Send(pkt)
}

// AuthAccount Packet
func AuthAccount(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified) {
		return
	}

	// skip 2 bytes
	reader.ReadUint16()

	// read and decrypt RSA block
	loginData := reader.ReadBytes(server.RSA_LOGIN_LENGTH)
	data, err := session.ServerInstance.RSA.Decrypt(loginData[:])
	if err != nil {
		log.Errorf("%s; Src: %s", err.Error(), session.GetEndPnt())
		session.Close()
		return
	}

	// extract name and pass
	name := string(bytes.Trim(data[:32], "\x00"))
	pass := string(bytes.Trim(data[32:], "\x00"))

	req := account.AuthRequest{UserId: name, Password: pass}
	res := account.AuthResponse{Status: account.None}
	err = session.RPC.Call(rpc.AuthCheck, &req, &res)

	// if server is down...
	if err != nil {
		res.Status = account.OutOfService
	}

	pkt := network.NewWriter(CSCAuthAccount)
	pkt.WriteByte(res.Status)
	pkt.WriteInt32(res.Id)
	pkt.WriteInt16(0x00)
	pkt.WriteByte(len(res.CharList)) // server count
	pkt.WriteInt64(0x00)
	pkt.WriteInt32(0x00) // premium service id
	pkt.WriteInt32(0x00) // premium service expire date
	pkt.WriteByte(0x00)
	pkt.WriteByte(res.SubPassChar) // subpassword exists for character
	pkt.WriteBytes(make([]byte, 7))
	pkt.WriteInt32(0x00) // language
	pkt.WriteString(res.AuthKey + "\x00")

	for _, value := range res.CharList {
		pkt.WriteByte(value.Server)
		pkt.WriteByte(value.Count)
	}

	session.Send(pkt)

	if res.Status == account.Normal {
		log.Infof("User `%s` successfully logged in.", name)

		session.Account = res.Id
		event.Trigger(event.PlayerLogin, session, name, true)

		// send url's
		URLToClient(session)

		// send normal system message
		session.Send(SystemMessage(message.Normal, 0))

		// create new periodic task to send server list periodically
		task := network.NewPeriodicTask(time.Second*5, func() {
			session.Send(ServerSate(session))
		})

		session.AddJob("ServerState", task)
	} else if res.Status == account.Online {
		session.Account = res.Id
		log.Infof("User `%s` double login attempt.", name)
	} else {
		log.Infof("User `%s` failed to log in.", name)
		event.Trigger(event.PlayerLogin, session, name, false)
	}
}
