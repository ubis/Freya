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

	serverStatus := 0 // 0 = normal, 1 = maintenance, 2 = outofservice
	pkt := network.NewWriter(CSCAuthAccount)
	pkt.WriteUint32(uint32(serverStatus))
	pkt.WriteUint32(0xFEFFFFFF)
	pkt.WriteUint32(0x02000000)
	session.Send(pkt)

	if serverStatus == 0 {
		conf := session.ServerConfig
		now := time.Now()
		pkt := network.NewWriter(NFYDisconnectTimer)
		pkt.WriteInt64(now.Unix() + int64(conf.AutoDisconnectTime))
		pkt.WriteByte(0)
		session.Send(pkt)
	}
}

// Authenticate Packet
func Authenticate(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified) {
		return
	}

	// 4 bytes padding
	reader.ReadUint32()

	// read and decrypt RSA block
	loginData := reader.ReadBytes(server.RSA_LOGIN_LENGTH)
	data, err := session.ServerInstance.RSA.Decrypt(loginData[:])
	if err != nil {
		log.Errorf("%s; Src: %s", err.Error(), session.GetEndPnt())
		session.Close()
		return
	}

	// extract name and pass
	name := string(bytes.Trim(data[:129], "\x00"))
	pass := string(bytes.Trim(data[129:], "\x00"))

	req := account.AuthRequest{UserId: name, Password: pass}
	res := account.AuthResponse{Status: account.None}
	err = session.RPC.Call(rpc.AuthCheck, &req, &res)

	// if server is down...
	if err != nil {
		res.Status = account.OutOfService
	}

	pkt := network.NewWriter(CSCAuthenticate)
	pkt.WriteBool(true) // keep alive
	pkt.WriteUint32(0)  // unknown
	pkt.WriteUint32(0)  // unknown

	// login status
	if res.Status == account.Normal || res.Status == account.Online {
		pkt.WriteInt32(1)
	} else {
		pkt.WriteInt32(0)
	}

	pkt.WriteUint32(0)         // not extended
	pkt.WriteInt32(res.Status) // account status
	session.Send(pkt)

	if res.Status != account.Normal {
		log.Infof("User `%s` failed to log in.", name)
		event.Trigger(event.PlayerLogin, session, name, false)
		return
	}

	conf := session.ServerConfig

	pkt = network.NewWriter(NFYAuthTimer)
	pkt.WriteUint32(conf.AutoDisconnectTime)
	session.Send(pkt)

	URLToClient(session)

	pkt = network.NewWriter(CSCAuthenticate)
	pkt.WriteBool(true) // keep alive
	pkt.WriteUint32(0)  // unknown
	pkt.WriteUint32(0)  // unknown

	// TODO: Migrate account.Online case

	// login status
	if res.Status == account.Normal || res.Status == account.Online {
		pkt.WriteInt32(1)
	} else {
		pkt.WriteInt32(0)
	}

	pkt.WriteUint32(0x11)     // extended
	pkt.WriteByte(res.Status) // account status
	pkt.WriteBytes(make([]byte, 55))
	pkt.WriteBytes(make([]byte, 32)) // TODO: This is a 32 byte auth key with following charset: 0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz
	pkt.WriteBytes(make([]byte, 3))

	for _, value := range res.CharList {
		pkt.WriteByte(value.Server)
		pkt.WriteByte(value.Count)
	}

	maxCharCount := 128
	pkt.WriteBytes(make([]byte, maxCharCount-len(res.CharList))) // max char count
	session.Send(pkt)

	pkt = network.NewWriter(NFYAuthTimer)
	pkt.WriteUint32(conf.AutoDisconnectTime)
	session.Send(pkt)

	log.Infof("User `%s` successfully logged in.", name)

	session.Account = res.Id
	event.Trigger(event.PlayerLogin, session, name, true)

	// send normal system message
	session.Send(SystemMessage(message.Normal, 0))

	// create new periodic task to send server list periodically
	task := network.NewPeriodicTask(time.Second*5, func() {
		session.Send(ServerSate(session))
	})
	session.AddJob("ServerState", task)
}
