package packet

import (
    "share/network"
    "share/models/account"
    "share/rpc"
)

// GetMyChartr Packet
func GetMyChartr(session *network.Session, reader *network.Reader) {
    if !session.Data.Verified {
        log.Error("Unauthorized connection from", session.GetEndPnt())
        session.Close()
        return
    }

    // fetch subpassword
    var req = account.SubPasswordReq{session.Data.AccountId}
    var res = account.SubPassword{}
    g_RPCHandler.Call(rpc.FetchSubPassword, req, &res)

    session.Data.SubPassword = &res

    var subpasswdExist = 0
    if res.Password != "" {
        subpasswdExist = 1
    }

    var packet = network.NewWriter(GETMYCHARTR)
    packet.WriteInt32(subpasswdExist)
    packet.WriteBytes(make([]byte, 9))
    packet.WriteInt32(0x00) // selected character id
    packet.WriteInt32(0x00) // slot order

    session.Send(packet)
}