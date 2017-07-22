package packet

import (
    "net"
    "encoding/binary"
    "share/network"
    "share/rpc"
    "share/models/server"
    "share/models/account"
)

// PreServerEnvRequest Packet
func PreServerEnvRequest(session *network.Session, reader *network.Reader) {
    var packet = network.NewWriter(PRE_SERVER_ENV_REQUEST)

    packet.WriteBytes(make([]byte, 4113))

    session.Send(packet)
}

// URLToClient Packet which is NFY
func URLToClient(session *network.Session) {
    var cash_url        = g_ServerConfig.CashWeb_URL
    var cash_odc_url    = g_ServerConfig.CashWeb_Odc_URL
    var cash_charge_url = g_ServerConfig.CashWeb_Charge_URL
    var guildweb_url    = g_ServerConfig.GuildWeb_URL
    var sns_url         = g_ServerConfig.Sns_URL

    var dataLen = len(cash_url) + 4
    dataLen += len(cash_odc_url) + 4
    dataLen += len(cash_charge_url) + 4
    dataLen += len(guildweb_url) + 4
    dataLen += len(sns_url) + 4

    var packet = network.NewWriter(URLTOCLIENT)
    packet.WriteInt16(dataLen + 2)
    packet.WriteInt16(dataLen)

    packet.WriteInt32(len(cash_url))
    packet.WriteString(cash_url)
    packet.WriteInt32(len(cash_odc_url))
    packet.WriteString(cash_odc_url)
    packet.WriteInt32(len(cash_charge_url))
    packet.WriteString(cash_charge_url)
    packet.WriteInt32(len(guildweb_url))
    packet.WriteString(guildweb_url)
    packet.WriteInt32(len(sns_url))
    packet.WriteString(sns_url)

    session.Send(packet)
}

// SystemMessg Packet which is NFY
func SystemMessg(session *network.Session, message byte) {
    var packet = network.NewWriter(SYSTEMMESSG)
    packet.WriteByte(message)
    packet.WriteByte(0x00)  // DataLength
    packet.WriteByte(0x00)  // Data
    session.Send(packet)
}

// ServerState Packet which is NFY
func ServerSate(session *network.Session) {
    var serverList = server.SvrListResponse{}
    g_RPCHandler.Call(rpc.ServerList, server.SvrListRequest{}, &serverList)

    var svr = serverList.Servers

    var packet = network.NewWriter(SERVERSTATE)
    packet.WriteByte(len(svr))

    for i := 0; i < len(svr); i ++ {
        packet.WriteByte(svr[i].Id)
        packet.WriteByte(svr[i].Hot)  // 0x10 = HOT! Flag; or bit_set(5)
        packet.WriteInt32(0x00)
        packet.WriteByte(len(svr[i].Channels))

        for j := 0; j < len(svr[i].Channels); j ++ {
            var channel = svr[i].Channels[j]
            var ip = binary.LittleEndian.Uint32(net.ParseIP(channel.Ip)[12:16])

            packet.WriteByte(channel.Id);
            packet.WriteUint16(channel.CurrentUsers);
            packet.WriteUint16(0x00);
            packet.WriteUint16(0xFFFF);
            packet.WriteUint16(0x00);
            packet.WriteUint16(0x00);
            packet.WriteUint32(0x00);
            packet.WriteUint16(0x00);
            packet.WriteUint16(0x00);
            packet.WriteUint16(0x00);
            packet.WriteByte(0x00);
            packet.WriteByte(0x00);
            packet.WriteByte(0x00);
            packet.WriteByte(0xFF);
            packet.WriteUint16(channel.MaxUsers);
            packet.WriteUint32(ip);
            packet.WriteUint16(channel.Port);
            packet.WriteUint32(channel.Type);
        }
    }

    session.Send(packet)
}

// VerifyLinks
func VerifyLinks(session *network.Session, reader *network.Reader) {
    var timestamp = reader.ReadUint32();
    var count     = reader.ReadUint16();
    var channel   = reader.ReadByte();
    var server    = reader.ReadByte();
    var magickey  = reader.ReadInt32();

    if magickey != int32(g_ServerConfig.MagicKey) {
        log.Errorf(
            "Invalid Client MagicKey (Required: %d, detected: %d, id: %d, src: %s",
        g_ServerConfig.MagicKey, magickey, session.Data.AccountId, session.GetEndPnt())
        return
    }

    var send = account.UserVerify{timestamp, count, server, channel, session.Data.AccountId}
    var recv = account.UserVerifyRecv{}

    g_RPCHandler.Call(rpc.UserVerify, send, &recv)

    var packet = network.NewWriter(VERIFYLINKS)
    packet.WriteByte(server)
    packet.WriteByte(channel)

    if recv.Verified {
        packet.WriteByte(0x01)
    } else{
        packet.WriteByte(0x00)
    }

    session.Send(packet)
}