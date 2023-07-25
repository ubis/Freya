package packet

import (
	"bytes"
	"time"

	"github.com/ubis/Freya/cmd/loginserver/rsa"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/account"
	"github.com/ubis/Freya/share/models/message"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

// PublicKey Packet
func PublicKey(session *network.Session, reader *network.Reader) {
	var rsa = g_ServerSettings.RSA
	var key = rsa.PublicKey

	var packet = network.NewWriter(PUBLIC_KEY)
	packet.WriteByte(0x01)
	packet.WriteUint16(len(key))
	packet.WriteBytes(key[:])

	session.Send(packet)
}

// AuthAccount Packet
func AuthAccount(session *network.Session, reader *network.Reader) {
	if session.Data.Verified != true {
		log.Errorf("Session version is not verified! Src: %s", session.GetEndPnt())
		session.Close()
		return
	}

	// skip 2 bytes
	reader.ReadUint16()

	// read and decrypt RSA block
	var loginData = reader.ReadBytes(rsa.RSA_LOGIN_LENGTH)
	var data, err = g_ServerSettings.RSA.Decrypt(loginData[:])
	if err != nil {
		log.Errorf("%s; Src: %s", err.Error(), session.GetEndPnt())
		session.Close()
		return
	}

	// extract name and pass
	var name = string(bytes.Trim(data[:32], "\x00"))
	var pass = string(bytes.Trim(data[32:], "\x00"))

	var r = account.AuthResponse{Status: account.None}
	err = g_RPCHandler.Call(rpc.AuthCheck, account.AuthRequest{name, pass}, &r)

	// if server is down...
	if err != nil {
		r.Status = account.OutOfService
	}

	var packet = network.NewWriter(AUTHACCOUNT)
	packet.WriteByte(r.Status)
	packet.WriteInt32(r.Id)
	packet.WriteInt16(0x00)
	packet.WriteByte(len(r.CharList)) // server count
	packet.WriteInt64(0x00)
	packet.WriteInt32(0x00) // premium service id
	packet.WriteInt32(0x00) // premium service expire date
	packet.WriteByte(0x00)
	packet.WriteByte(r.SubPassChar) // subpassword exists for character
	packet.WriteBytes(make([]byte, 7))
	packet.WriteInt32(0x00) // language
	packet.WriteString(r.AuthKey + "\x00")

	for _, value := range r.CharList {
		packet.WriteByte(value.Server)
		packet.WriteByte(value.Count)
	}

	session.Send(packet)

	if r.Status == account.Normal {
		log.Infof("User `%s` successfully logged in.", name)

		session.Data.AccountId = r.Id
		session.Data.LoggedIn = true

		// send url's
		URLToClient(session)

		// send normal system message
		session.Send(SystemMessg(message.Normal, 0))

		// send server list periodically
		var t = time.NewTicker(time.Second * 5)
		go func(s *network.Session) {
			for {
				if !s.Connected {
					break
				}

				s.Send(ServerSate())
				<-t.C
			}
		}(session)
	} else if r.Status == account.Online {
		session.Data.AccountId = r.Id
		log.Infof("User `%s` double login attempt.", name)
	} else {
		log.Infof("User `%s` failed to log in.", name)
	}
}
