package packet

import (
    "share/network"
    "share/rpc"
    "share/models/subpasswd"
)

// GetMyChartr Packet
func GetMyChartr(session *network.Session, reader *network.Reader) {
    if !session.Data.Verified {
        log.Error("Unauthorized connection from", session.GetEndPnt())
        session.Close()
        return
    }

    // fetch subpassword
    var req = subpasswd.FetchReq{session.Data.AccountId}
    var res = subpasswd.FetchRes{}
    g_RPCHandler.Call(rpc.FetchSubPassword, req, &res)

    session.Data.SubPassword = &res.Details

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