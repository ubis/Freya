package packet

import (
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/character"
	"github.com/ubis/Freya/share/models/inventory"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

func NotifyStorageExchange(session *network.Session, result byte) {
	packet := network.NewWriter(STORAGE_EXCHANGE_MOVE)
	packet.WriteByte(result)
	packet.WriteInt32(0)

	session.Send(packet)
}

func UpdatePlayerStorage(session *network.Session, packet string, deleteSlot uint16, createSlot uint16) (inventory.Item, error) {
	id, err := getSessionCharId(session)
	if err != nil {
		return inventory.Item{}, err
	}

	send := character.ItemMoveReq{Server: byte(g_ServerSettings.ServerId), Id: id, DeleteSlot: deleteSlot, CreateSlot: createSlot}
	receive := character.ItemMoveRes{}
	err = g_RPCHandler.Call(packet, send, &receive)
	if err != nil {
		return inventory.Item{}, err
	}

	if receive.Result != nil {
		return inventory.Item{}, receive.Result
	}

	ctx, err := parseSessionContext(session)
	if err != nil {
		return inventory.Item{}, err
	}

	ctx.mutex.Lock()
	ctx.char.Equipment.Remove(deleteSlot)
	ctx.char.Inventory.Set(createSlot, receive.Item)
	ctx.mutex.Unlock()

	return receive.Item, nil
}

func StorageExchangeMove(session *network.Session, reader *network.Reader) {
	isEquip := reader.ReadUint32() == 1
	deleteSlot := uint16(reader.ReadUint32())
	isInventory := reader.ReadUint32() == 1
	createSlot := uint16(reader.ReadUint32())

	var err error = nil

	if isEquip && !isInventory {
		_, updateErr := UpdatePlayerStorage(session, rpc.UnEquipItem, deleteSlot, createSlot)
		if updateErr != nil {
			log.Error(err)
			NotifyStorageExchange(session, 0x00)
			return
		}

		id, err2 := getSessionCharId(session)
		if err2 != nil {
			log.Error(err2.Error())
			NotifyStorageExchange(session, 0x00)
			return
		}

		pkt := network.NewWriter(NFY_ITEM_UNEQUIP)
		pkt.WriteInt32(id)
		pkt.WriteInt16(deleteSlot)
		g_NetworkManager.SendToAll(pkt)

		NotifyStorageExchange(session, 0x01)

		return
	}

	if isInventory && !isEquip {
		item, updateErr := UpdatePlayerStorage(session, rpc.EquipItem, deleteSlot, createSlot)
		if updateErr != nil {
			log.Error(updateErr)
			NotifyStorageExchange(session, 0x00)
			return
		}

		id, charIdErr := getSessionCharId(session)
		if charIdErr != nil {
			log.Error(charIdErr.Error())
			NotifyStorageExchange(session, 0x00)
			return
		}

		pkt := network.NewWriter(NFY_ITEM_EQUIP)
		pkt.WriteInt32(id)
		pkt.WriteInt32(item.Kind)
		pkt.WriteInt16(item.Slot)
		pkt.WriteByte(0x00)
		pkt.WriteByte(0x00)
		pkt.WriteByte(0x00)
		pkt.WriteByte(0x00)
		pkt.WriteByte(0x00)

		g_NetworkManager.SendToAll(pkt)

		NotifyStorageExchange(session, 0x01)
		return
	}

	if isEquip && isInventory {
		_, updateErr := UpdatePlayerStorage(session, rpc.ChangeEquipItemSlot, deleteSlot, createSlot)
		if updateErr != nil {
			log.Error(updateErr)
			NotifyStorageExchange(session, 0x00)
			return

			//TODO Check and add notify packets
		}

		NotifyStorageExchange(session, 0x01)
		return
	}

	if !isEquip && !isInventory {
		_, updateErr := UpdatePlayerStorage(session, rpc.ChangeInventoryItemSlot, uint16(deleteSlot), uint16(createSlot))

		if updateErr != nil {
			log.Error(updateErr)
			NotifyStorageExchange(session, 0x00)
			return
		}

		NotifyStorageExchange(session, 0x01)
		return
	}

	NotifyStorageExchange(session, 0x00)
}
