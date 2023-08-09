package packet

import (
	"github.com/ubis/Freya/share/models/character"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
	"log"
)

func MoveItem(session *network.Session, reader *network.Reader) {
	//EQUIP -> INVENTORY
	//01 00 00 00 02 00 00 00 00 00 00 00 00 00 00 00 00 00 F4 DC EB FF FF FF 00 00 00 00
	//01 00 00 00 04 00 00 00 00 00 00 00 02 00 00 00 00 24 48 E1 EB FF FF FF 00 00 00 00

	//INVENTORY -> EQUIP
	//00 00 00 00 02 00 00 00 01 00 00 00 02 00 00 00 00 00 00 00 EB FF FF FF 00 00 00 00
	//00 00 00 00 06 00 00 00 01 00 00 00 01 00 00 00 00 00 00 00 EB FF FF FF 00 00 00 00

	//INVENTORY -> INVENTORY
	//00 00 00 00 04 00 00 00 00 00 00 00 00 00 00 00 00 0E 02 00 EB FF FF FF 00 00 00 00
	//00 00 00 00 00 00 00 00 00 00 00 00 01 00 00 00 00 0E 03 00 EB FF FF FF 00 00 00 00

	//EQUIP -> EQUIP
	//01 00 00 00 08 00 00 00 01 00 00 00 09 00 00 00 00 00 18 0B EB FF FF FF 00 00 00 00
	//01 00 00 00 09 00 00 00 01 00 00 00 12 00 00 00 00 00 18 0B EB FF FF FF 00 00 00 00

	//log.Printf("%X", reader.ReadBytes(28))

	equip := reader.ReadUint32()
	deleteSlot := reader.ReadUint32()
	inventory := reader.ReadInt32()
	createSlot := reader.ReadUint32()

	var packet = network.NewWriter(MOVE_ITEM)

	packet.WriteByte(0x01)
	packet.WriteInt64(0)

	session.Send(packet)

	id, err := getSessionCharId(session)

	if err != nil {
		log.Printf(err.Error())
		return
	}

	if equip > 0 {
		if inventory > 0 {
			var send = character.ItemEquipReq{byte(g_ServerSettings.ServerId), id, uint16(deleteSlot), uint16(createSlot)}
			var recv = character.ItemEquipRes{}
			g_RPCHandler.Call(rpc.MoveItemEquToEqu, send, &recv)
		} else {
			var send = character.ItemEquipReq{byte(g_ServerSettings.ServerId), id, uint16(deleteSlot), uint16(createSlot)}
			var recv = character.ItemEquipRes{}
			g_RPCHandler.Call(rpc.MoveItemEquToInv, send, &recv)
		}
	} else {
		if inventory > 0 {
			var send = character.ItemEquipReq{byte(g_ServerSettings.ServerId), id, uint16(deleteSlot), uint16(createSlot)}
			var recv = character.ItemEquipRes{}
			g_RPCHandler.Call(rpc.MoveItemInvToEqu, send, &recv)
		} else {
			var send = character.ItemEquipReq{byte(g_ServerSettings.ServerId), id, uint16(deleteSlot), uint16(createSlot)}
			var recv = character.ItemEquipRes{}
			g_RPCHandler.Call(rpc.MoveItemInvToInv, send, &recv)
		}
	}
}
