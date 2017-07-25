package packet

import (
    "time"
    "bytes"
    "share/network"
    "share/rpc"
    "share/models/subpasswd"
    "golang.org/x/crypto/bcrypt"
)

// SubPasswordSet Packet
func SubPasswordSet(session *network.Session, reader *network.Reader) {
    var passwd = string(bytes.Trim(reader.ReadBytes(10), "\x00"))
    reader.ReadBytes(5)

    var question = reader.ReadInt32()
    var answer   = string(bytes.Trim(reader.ReadBytes(16), "\x00"))

    var packet = network.NewWriter(SUBPW_SET)

    var sub      = session.Data.SubPassword
    var verified = &sub.Verified
    if sub.Password == "" {
        // verified because user is creating for the first time
        *verified = true
    }

    if len(passwd) < 4 || question < 1 || question > 10 || !*verified  {
        packet.WriteInt32(0x00) // failed
        packet.WriteInt32(0x00)
        packet.WriteInt32(0x01)
        packet.WriteInt32(0x00)

        session.Send(packet)
        return
    }

    if sub.Password == "" {
        // creating sub password for the first time
        // check answer
        if len(answer) < 4 {
            packet.WriteInt32(0x00) // failed
            packet.WriteInt32(0x00)
            packet.WriteInt32(0x01)
            packet.WriteInt32(0x00)

            session.Send(packet)
            return
        }

        var hash, _  = bcrypt.GenerateFromPassword([]byte(answer), bcrypt.DefaultCost)
        sub.Answer   = string(hash)
        sub.Question = byte(question)
        sub.Expires  = time.Now()
    }

    var hash, _  = bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
    sub.Password = string(hash)

    // update to db
    var req = subpasswd.SetReq{session.Data.AccountId, *sub}
    var res = subpasswd.SetRes{}
    var err = g_RPCHandler.Call(rpc.SetSubPassword, req, &res)

    if err == nil && res.Success {
        packet.WriteInt32(0x01) // success
    } else {
        packet.WriteInt32(0x00) // failed
    }

    *verified = false

    packet.WriteInt32(0x00)
    packet.WriteInt32(0x01)
    packet.WriteInt32(0x00)

    session.Send(packet)
}

// SubPasswordCheckRequest Packet
func SubPasswordCheckRequest(session *network.Session, reader *network.Reader) {
    var sub  = session.Data.SubPassword
    var left = sub.Expires.Sub(time.Now())

    var packet = network.NewWriter(SUBPW_CHECK_REQ)

    if sub.Password == "" {
        // need to create first
        packet.WriteInt32(0x00)
    } else {
        // now check actual time
        if left > 0 {
            // no verification needed
            packet.WriteInt32(0x00)
        } else {
            // verification is needed
            packet.WriteInt32(0x01)
        }
    }

    session.Send(packet)
}

// SubPasswordCheck Packet
func SubPasswordCheck(session *network.Session, reader *network.Reader) {
    var password = string(bytes.Trim(reader.ReadBytes(10), "\x00"))
    reader.ReadBytes(5)
    var hours = reader.ReadInt32()

    var sub = session.Data.SubPassword
    var err = bcrypt.CompareHashAndPassword([]byte(sub.Password), []byte(password))

    var packet = network.NewWriter(SUBPW_CHECK)

    if hours < 0 || hours > 4 {
        packet.WriteInt32(0x00) // failed
        packet.WriteByte(sub.FailTimes)
        packet.WriteInt32(0x00)
        packet.WriteInt32(0x01)

        sub.FailTimes ++
        session.Send(packet)
        return
    }

    if err != nil {
        packet.WriteInt32(0x00) // failed
        sub.FailTimes ++
    } else {
        sub.Expires = sub.Expires.Add(time.Hour * time.Duration(hours))
        var req = subpasswd.SetReq{session.Data.AccountId, *sub}
        var res = subpasswd.SetRes{}
        err = g_RPCHandler.Call(rpc.SetSubPassword, req, &res)

        if err != nil || !res.Success {
            packet.WriteInt32(0x00) // failed
            sub.FailTimes ++
        } else {
            packet.WriteInt32(0x01) // success
            sub.FailTimes = 0
            sub.Verified  = true
        }
    }

    packet.WriteByte(sub.FailTimes)
    packet.WriteInt32(0x00)
    packet.WriteInt32(0x01)

    session.Send(packet)
}

// SubPasswordFindRequest Packet
func SubPasswordFindRequest(session *network.Session, reader *network.Reader) {
    var sub = session.Data.SubPassword

    var packet = network.NewWriter(SUBPW_FIND_REQ)
    packet.WriteInt32(sub.Question)
    packet.WriteInt32(sub.Question)
    packet.WriteInt32(0x01)

    session.Send(packet)
}

// SubPasswordFind Packet
func SubPasswordFind(session *network.Session, reader *network.Reader) {
    reader.ReadBytes(8)
    var answer = string(bytes.Trim(reader.ReadBytes(16), "\x00"))

    var sub = session.Data.SubPassword
    var err = bcrypt.CompareHashAndPassword([]byte(sub.Answer), []byte(answer))

    var packet = network.NewWriter(SUBPW_FIND)

    if err != nil {
        packet.WriteInt32(0x00) // failed
        sub.FailTimes ++
    } else {
        packet.WriteInt32(0x01) // success
        sub.FailTimes = 0
        sub.Verified  = true
    }

    packet.WriteByte(sub.FailTimes)
    packet.WriteInt32(0x01)

    session.Send(packet)
}

// SubPasswordDelRequest Packet
func SubPasswordDelRequest(session *network.Session, reader *network.Reader) {
    reader.ReadBytes(4)
    var password = string(bytes.Trim(reader.ReadBytes(10), "\x00"))

    var sub = session.Data.SubPassword
    var err = bcrypt.CompareHashAndPassword([]byte(sub.Password), []byte(password))

    var packet = network.NewWriter(SUBPW_DEL_REQ)

    if err != nil {
        packet.WriteInt32(0x00) // failed
        sub.FailTimes ++
    } else {
        packet.WriteInt32(0x01) // success
        sub.FailTimes = 0
        sub.Verified  = true
    }

    packet.WriteByte(sub.FailTimes)
    packet.WriteInt32(0x01)

    session.Send(packet)
}

// SubPasswordDel Packet
func SubPasswordDel(session *network.Session, reader *network.Reader) {
    var packet = network.NewWriter(SUBPW_DEL)

    var sub = session.Data.SubPassword

    if !sub.Verified  {
        packet.WriteInt32(0x00) // failed
        packet.WriteInt32(0x01)

        session.Send(packet)
        return
    }

    // update to db
    var req = subpasswd.SetReq{Account: session.Data.AccountId}
    var res = subpasswd.SetRes{}
    var err = g_RPCHandler.Call(rpc.RemoveSubPassword, req, &res)

    if err == nil && res.Success {
        *sub = subpasswd.Details{}
        packet.WriteInt32(0x01) // success
    } else {
        packet.WriteInt32(0x00) // failed
    }

    packet.WriteInt32(0x01)

    session.Send(packet)
}

// SubPasswordChangeQARequest Packet
func SubPasswordChangeQARequest(session *network.Session, reader *network.Reader) {
    reader.ReadBytes(4)
    var password = string(bytes.Trim(reader.ReadBytes(10), "\x00"))

    var sub = session.Data.SubPassword
    var err = bcrypt.CompareHashAndPassword([]byte(sub.Password), []byte(password))

    var packet = network.NewWriter(SUBPW_CHG_QA_REQ)

    if err != nil {
        packet.WriteInt32(0x00) // failed
        sub.FailTimes ++
    } else {
        packet.WriteInt32(0x01) // success
        sub.FailTimes = 0
        sub.Verified  = true
    }

    packet.WriteByte(sub.FailTimes)
    packet.WriteInt32(0x01)

    session.Send(packet)
}

// SubPasswordChangeQA Packet
func SubPasswordChangeQA(session *network.Session, reader *network.Reader) {
    reader.ReadBytes(4)
    var question = reader.ReadInt32()
    var answer   = string(bytes.Trim(reader.ReadBytes(16), "\x00"))

    var packet = network.NewWriter(SUBPW_CHG_QA)

    var sub = session.Data.SubPassword

    if len(answer) < 4 || question < 1 || question > 10 || !sub.Verified  {
        packet.WriteInt32(0x00) // failed
        packet.WriteInt32(0x01)

        session.Send(packet)
        return
    }

    var hash, _  = bcrypt.GenerateFromPassword([]byte(answer), bcrypt.DefaultCost)
    sub.Answer   = string(hash)
    sub.Question = byte(question)

    // update to db
    var req = subpasswd.SetReq{session.Data.AccountId, *sub}
    var res = subpasswd.SetRes{}
    var err = g_RPCHandler.Call(rpc.SetSubPassword, req, &res)

    if err == nil && res.Success {
        packet.WriteInt32(0x01) // success
    } else {
        packet.WriteInt32(0x00) // failed
    }

    sub.Verified = false
    packet.WriteInt32(0x01)

    session.Send(packet)
}

// CharacterDeleteCheckSubPassword Packet
func CharacterDeleteCheckSubPassword(session *network.Session, reader *network.Reader) {
    var password = string(bytes.Trim(reader.ReadBytes(10), "\x00"))

    var sub = session.Data.SubPassword
    var err = bcrypt.CompareHashAndPassword([]byte(sub.Password), []byte(password))

    var packet = network.NewWriter(CHAR_DEL_CHK_SUBPW)

    if err != nil {
        packet.WriteInt32(0x00) // failed
        sub.FailTimes ++
    } else {
        packet.WriteInt32(0x01) // success
        sub.FailTimes = 0
        sub.Verified  = true
    }

    packet.WriteByte(sub.FailTimes)

    session.Send(packet)
}