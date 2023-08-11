package packet

import (
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/character"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

func SendStorageExchangeMoveStatusPacket(session *network.Session, result byte) {
	packet := network.NewWriter(STORAGE_EXCHANGE_MOVE)
	packet.WriteByte(result)
	packet.WriteInt32(0)

	session.Send(packet)
}

func SendStorageExchangeMoveMasterServerPacket(packet string, id int32, deleteSlot uint16, createSlot uint16) character.ItemMoveRes {
	send := character.ItemMoveReq{Server: byte(g_ServerSettings.ServerId), Id: id, DeleteSlot: deleteSlot, CreateSlot: createSlot}
	receive := character.ItemMoveRes{}
	err := g_RPCHandler.Call(packet, send, &receive)
	if err != nil {
		log.Error(err.Error())
	}

	return receive
}

func StorageExchangeMove(session *network.Session, reader *network.Reader) {
	equip := reader.ReadUint32()
	deleteSlot := uint16(reader.ReadUint32())
	inventory := reader.ReadUint32()
	createSlot := uint16(reader.ReadUint32())

	ctx, err := parseSessionContext(session)
	if err != nil {
		log.Error(err.Error())
		SendStorageExchangeMoveStatusPacket(session, 0x00)
		return
	}
	playerChar := ctx.char

	isEquip := equip == 1
	isInventory := inventory == 1

	if isEquip && !isInventory {
		result := SendStorageExchangeMoveMasterServerPacket(rpc.UnEquipItem, playerChar.Id, deleteSlot, createSlot)

		if result.Result == nil {
			playerChar.Equipment.Remove(deleteSlot)
			//playerChar.Inventory.Set(createSlot, result.Item)

			pkt := network.NewWriter(NFY_ITEM_UNEQUIP)
			pkt.WriteInt32(playerChar.Id)
			pkt.WriteInt16(deleteSlot)

			g_NetworkManager.SendToAll(pkt)

			SendStorageExchangeMoveStatusPacket(session, 0x01)
			return
		} else {
			log.Error(result.Result)
			SendStorageExchangeMoveStatusPacket(session, 0x00)
			return
		}
	}

	if isInventory && !isEquip {
		result := SendStorageExchangeMoveMasterServerPacket(rpc.EquipItem, playerChar.Id, deleteSlot, createSlot)

		if result.Result == nil {
			playerChar.Inventory.Remove(deleteSlot)
			playerChar.Equipment.Set(createSlot, result.Item)

			pkt := network.NewWriter(NFY_ITEM_EQUIP)
			pkt.WriteInt32(playerChar.Id)
			pkt.WriteInt32(result.Item.Kind)
			pkt.WriteInt16(result.Item.Slot)
			pkt.WriteByte(0x00)
			pkt.WriteByte(0x00)
			pkt.WriteByte(0x00)
			pkt.WriteByte(0x00)
			pkt.WriteByte(0x00)

			g_NetworkManager.SendToAll(pkt)

			SendStorageExchangeMoveStatusPacket(session, 0x01)
			return
		} else {
			log.Error(result.Result)
			SendStorageExchangeMoveStatusPacket(session, 0x00)
			return
		}
	}

	if isEquip && isInventory {
		result := SendStorageExchangeMoveMasterServerPacket(rpc.ChangeEquipItemSlot, playerChar.Id, deleteSlot, createSlot)
		if result.Result == nil {
			playerChar.Equipment.Remove(deleteSlot)
			playerChar.Equipment.Set(createSlot, result.Item)

			//TODO Check and add notify packets

			SendStorageExchangeMoveStatusPacket(session, 0x01)
			return
		} else {
			log.Error(result.Result)
			SendStorageExchangeMoveStatusPacket(session, 0x00)
			return
		}
	}

	if !isEquip && !isInventory {
		result := SendStorageExchangeMoveMasterServerPacket(rpc.ChangeInventoryItemSlot, playerChar.Id, uint16(deleteSlot), uint16(createSlot))

		if result.Result == nil {
			playerChar.Inventory.Remove(deleteSlot)
			//playerChar.Inventory.Set(createSlot, result.Item)

			SendStorageExchangeMoveStatusPacket(session, 0x01)
			return
		} else {
			log.Error(result.Result)
			SendStorageExchangeMoveStatusPacket(session, 0x00)
			return
		}
	}

	SendStorageExchangeMoveStatusPacket(session, 0x00)
}
