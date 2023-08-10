package rpc

import (
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/character"
	"github.com/ubis/Freya/share/models/inventory"
	"github.com/ubis/Freya/share/rpc"
)

func MoveItemInvToEqu(c *rpc.Client, r *character.ItemEquipReq, s *character.ItemEquipRes) error {
	var db = g_DatabaseManager.Get(r.Server)

	var inv = inventory.Inventory{}
	inv.Init()

	var equ = inventory.Equipment{}
	equ.Init()

	var item inventory.Item
	query := db.QueryRow("SELECT kind, serials, opt, slot, expire "+
		"FROM characters_inventory WHERE id = ? AND slot = ?", r.Id, r.DeleteSlot)
	err := query.Scan(&item.Kind, &item.Serials, &item.Option, &item.Slot, &item.Expire)

	//log.Error("EQUIP ITEM ", item)
	//log.Error("EQUIP [DATABASE]", err)
	//log.Debug("MOVE E TO E Request ", r)

	if err == nil {
		inv.Remove(item.Slot)
		var _, err3 = db.Queryx(
			"DELETE "+
				"FROM characters_inventory "+
				"WHERE id = ? AND slot = ?", r.Id, item.Slot)
		if err3 != nil {
			log.Error("[DATABASE]", err3)
		}
		item.Slot = r.CreateSlot
		equ.Set(item.Slot, item)
		db.MustExec("INSERT INTO characters_equipment "+
			"(id, kind, serials, opt, slot, expire) VALUES (?, ?, ?, ?, ?, ?)",
			r.Id, item.Kind, item.Serials, item.Option, item.Slot, item.Expire)
		s.ItemKind = uint32(item.Kind)
	}

	if err != nil {
		log.Error("[DATABASE]", err)
	}

	return nil
}

func MoveItemEquToInv(c *rpc.Client, r *character.ItemEquipReq, s *character.ItemEquipRes) error {
	var db = g_DatabaseManager.Get(r.Server)

	var inv = inventory.Inventory{}
	inv.Init()

	var equ = inventory.Equipment{}
	equ.Init()

	var item inventory.Item
	query := db.QueryRow("SELECT kind, serials, opt, slot, expire "+
		"FROM characters_equipment WHERE id = ? AND slot = ?", r.Id, r.DeleteSlot)
	err := query.Scan(&item.Kind, &item.Serials, &item.Option, &item.Slot, &item.Expire)

	//log.Error("UNEQUIP ITEM ", item)
	//log.Error("UNEQUIP [DATABASE]", err)
	//log.Debug("MOVE E TO I Request ", r)

	if err == nil {
		equ.Remove(item.Slot)
		var _, err3 = db.Queryx(
			"DELETE "+
				"FROM characters_equipment "+
				"WHERE id = ? AND slot = ?", r.Id, item.Slot)
		if err3 != nil {
			log.Error("[DATABASE]", err3)
		}
		item.Slot = r.CreateSlot
		inv.Set(item.Slot, item)
		db.MustExec("INSERT INTO characters_inventory "+
			"(id, kind, serials, opt, slot, expire) VALUES (?, ?, ?, ?, ?, ?)",
			r.Id, item.Kind, item.Serials, item.Option, item.Slot, item.Expire)
		s.ItemKind = uint32(item.Kind)
	}

	if err != nil {
		log.Error("[DATABASE]", err)
	}

	return nil
}

func MoveItemEquToEqu(c *rpc.Client, r *character.ItemEquipReq, s *character.ItemEquipRes) error {
	var db = g_DatabaseManager.Get(r.Server)

	var item inventory.Item
	query := db.QueryRow("SELECT kind, serials, opt, slot, expire "+
		"FROM characters_equipment WHERE id = ? AND slot = ?", r.Id, r.DeleteSlot)
	err := query.Scan(&item.Kind, &item.Serials, &item.Option, &item.Slot, &item.Expire)

	//log.Error("UNEQUIP ITEM ", item)
	//log.Error("UNEQUIP [DATABASE]", err)
	//log.Debug("MOVE E TO E Request ", r)

	if err == nil {
		db.MustExec("UPDATE characters_equipment SET slot = ? WHERE id = ? AND slot = ?", r.CreateSlot, r.Id, item.Slot)
		item.Slot = r.CreateSlot
		s.ItemKind = item.Kind
	}

	if err != nil {
		log.Error("[DATABASE]", err)
	}

	return nil
}

func MoveItemInvToInv(c *rpc.Client, r *character.ItemEquipReq, s *character.ItemEquipRes) error {
	var db = g_DatabaseManager.Get(r.Server)

	var item inventory.Item
	query := db.QueryRow("SELECT kind, serials, opt, slot, expire "+
		"FROM characters_inventory WHERE id = ? AND slot = ?", r.Id, r.DeleteSlot)
	err := query.Scan(&item.Kind, &item.Serials, &item.Option, &item.Slot, &item.Expire)

	//log.Error("UNEQUIP ITEM ", item)
	//log.Error("UNEQUIP [DATABASE]", err)
	//log.Debug("MOVE I TO I Request ", r)

	if err == nil {
		db.MustExec("UPDATE characters_inventory SET slot = ? WHERE id = ? AND slot = ?", r.CreateSlot, r.Id, item.Slot)
		item.Slot = r.CreateSlot
		s.ItemKind = item.Kind
	}

	if err != nil {
		log.Error("[DATABASE]", err)
	}

	return nil
}
