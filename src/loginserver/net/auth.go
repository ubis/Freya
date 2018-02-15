package net

import (
	"bytes"
	"loginserver/rsa"
	"share/log"
	"share/models/account"
	"share/models/message"
	"share/network"
	"share/rpc"
	"time"
)

// PublicKey Packet
func (p *Packet) PublicKey(s *network.Session, r *network.Reader) {
	var key = p.RSA.PublicKey

	var packet = network.NewWriter(PublicKey)
	packet.WriteByte(0x01)
	packet.WriteUint16(len(key))
	packet.WriteBytes(key[:])

	s.Send(packet)
}

// AuthAccount Packet
func (p *Packet) AuthAccount(s *network.Session, r *network.Reader) {
	if s.Data.Verified != true {
		log.Errorf("Session version is not verified! Src: %s", s.GetEndPnt())
		s.Close()
		return
	}

	// skip 2 bytes
	r.ReadUint16()

	// read and decrypt RSA block
	var loginData = r.ReadBytes(rsa.RSA_LOGIN_LENGTH)
	var data, err = p.RSA.Decrypt(loginData[:])
	if err != nil {
		log.Errorf("%s; Src: %s", err.Error(), s.GetEndPnt())
		s.Close()
		return
	}

	// extract name and pass
	var name = string(bytes.Trim(data[:32], "\x00"))
	var pass = string(bytes.Trim(data[32:], "\x00"))

	var req = account.AuthRequest{UserId: name, Password: pass}
	var rsp = account.AuthResponse{Status: account.None}
	err = p.RPC.Call(rpc.AuthCheck, req, &rsp)

	// if server is down...
	if err != nil {
		rsp.Status = account.OutOfService
	}

	var packet = network.NewWriter(AuthAccount)
	packet.WriteByte(rsp.Status)
	packet.WriteInt32(rsp.Id)
	packet.WriteInt16(0x00)
	packet.WriteByte(len(rsp.CharList)) // server count
	packet.WriteInt64(0x00)
	packet.WriteInt32(0x00) // premium service id
	packet.WriteInt32(0x00) // premium service expire date
	packet.WriteByte(0x00)
	packet.WriteByte(rsp.SubPassChar) // subpassword exists for character
	packet.WriteBytes(make([]byte, 7))
	packet.WriteInt32(0x00) // language
	packet.WriteString(rsp.AuthKey + "\x00")

	for _, value := range rsp.CharList {
		packet.WriteByte(value.Server)
		packet.WriteByte(value.Count)
	}

	s.Send(packet)

	if rsp.Status == account.Normal {
		log.Infof("User `%s` successfully logged in.", name)

		s.Data.AccountId = rsp.Id
		s.Data.LoggedIn = true

		// send url's
		p.URLToClient(s)

		// send normal system message
		s.Send(p.SystemMessg(message.Normal, 0))

		// send server list periodically
		var t = time.NewTicker(time.Second * 5)
		go func(s *network.Session, p *Packet) {
			for {
				if !s.Connected {
					break
				}

				s.Send(p.ServerState())
				<-t.C
			}
		}(s, p)
	} else if rsp.Status == account.Online {
		s.Data.AccountId = rsp.Id
		log.Infof("User `%s` double login attempt.", name)
	} else {
		log.Infof("User `%s` failed to log in.", name)
	}
}
