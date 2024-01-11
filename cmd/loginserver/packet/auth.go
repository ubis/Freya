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

	const (
		statusNormal int = iota
		statusMaintenance
		statusOutOfService
	)

	// here we could check if server is in the maintenance and/or whitelist mode
	// which could be set through IPC, through master server
	serverStatus := statusNormal

	pkt := network.NewWriter(CSCAuthAccount)
	pkt.WriteInt32(serverStatus)
	pkt.WriteUint32(0xFEFFFFFF)
	pkt.WriteUint32(0x02000000)

	session.Send(pkt)

	if serverStatus != statusNormal {
		return
	}

	// set-up disconnection timer
	// to-do: set this on server-side as well
	conf := session.ServerConfig
	session.Send(DisconnectTimer(conf.AutoDisconnectTime))
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
		res.Status = account.OutOfService2
	}

	keepAlive := res.Status == account.Normal || res.Status == account.Online

	pkt := network.NewWriter(CSCAuthenticate)
	pkt.WriteBool(keepAlive)   // keep alive
	pkt.WriteUint32(0)         // unknown
	pkt.WriteUint32(0)         // unknown
	pkt.WriteInt32(keepAlive)  // login status
	pkt.WriteUint32(0)         // not extended
	pkt.WriteInt32(res.Status) // account status

	session.Send(pkt)

	// set-up logging-in timer
	// to-do: set this on server-side as well
	conf := session.ServerConfig
	session.Send(AuthTimer(conf.AutoDisconnectTime))

	if !keepAlive {
		log.Infof("User `%s` failed to log in.", name)
		event.Trigger(event.PlayerLogin, session, name, false)
		return
	}

	// send URLs to the client
	URLToClient(session)

	pkt = network.NewWriter(CSCAuthenticate)
	pkt.WriteBool(keepAlive)  // keep alive
	pkt.WriteUint32(0)        // unknown
	pkt.WriteUint32(0)        // unknown
	pkt.WriteInt32(keepAlive) // login status
	pkt.WriteUint32(0x11)     // extended
	pkt.WriteByte(res.Status) // account status
	pkt.WriteBytes(make([]byte, 55))
	pkt.WriteString(res.AuthKey + "\x00")
	pkt.WriteBytes(make([]byte, 3))

	for _, value := range res.CharList {
		pkt.WriteByte(value.Server)
		pkt.WriteByte(value.Count)
	}

	maxCharCount := 128
	pkt.WriteBytes(make([]byte, maxCharCount-len(res.CharList))) // max char count
	session.Send(pkt)

	// set-up logging-in timer
	// to-do: set this on server-side as well
	session.Send(AuthTimer(conf.AutoDisconnectTime))

	log.Infof("User `%s` successfully logged in.", name)

	// set-up player account id
	session.Account = res.Id

	// trigger player login event
	event.Trigger(event.PlayerLogin, session, name, true)

	// send normal system message
	session.Send(SystemMessage(message.Normal, 0))

	// create new periodic task to send server list periodically
	task := network.NewPeriodicTask(time.Second*5, func() {
		session.Send(ServerSate(session))
	})
	session.AddJob("ServerState", task)
}
