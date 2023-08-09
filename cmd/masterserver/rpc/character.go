package rpc

import (
	"time"

	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/character"
	"github.com/ubis/Freya/share/models/inventory"
	"github.com/ubis/Freya/share/models/skills"
	"github.com/ubis/Freya/share/rpc"
)

// LoadCharacters RPC Call
func LoadCharacters(_ *rpc.Client, r *character.ListReq, s *character.ListRes) error {
	var db = g_DatabaseManager.Get(r.Server)
	var res = character.ListRes{List: make([]character.Character, 0, 6)}

	if db == nil {
		*s = res
		return nil
	}

	var rows, err = db.Queryx(
		"SELECT "+
			"id, name, level, world, x, y, alz, nation, sword_rank, magic_rank, "+
			"current_hp, max_hp, current_mp, max_mp, current_sp, max_sp, str_stat, "+
			"int_stat, dex_stat, pnt_stat, exp, war_exp, created "+
			"FROM characters "+
			"WHERE id >= ? AND id <= ?", r.Account*8, r.Account*8+5)

	if err != nil {
		log.Error("[DATABASE]", err)
		*s = res
		return nil
	}

	// iterate over each row
	for rows.Next() {
		var c = character.Character{}
		var err = rows.StructScan(&c)

		if err == nil {
			c.Style = LoadStyle(db, c.Id)
			c.Equipment = LoadEquipment(db, c.Id)

			res.List = append(res.List, c)
		} else {
			log.Error("[DATABASE]", err)
			*s = res
			return nil
		}
	}

	// load metadata
	db.Get(&res.SlotOrder,
		"SELECT slot_order FROM lobby_metadata WHERE id = ?", r.Account)
	db.Get(&res.LastId, "SELECT last_char FROM lobby_metadata WHERE id = ?", r.Account)

	*s = res
	return nil
}

// LoadStyle Database Call
func LoadStyle(db *sqlx.DB, id int32) character.Style {
	var style = character.Style{}
	var err = db.Get(&style,
		"SELECT battle_style, rank, face, color, hair, aura, gender, show_helmet "+
			"FROM characters "+
			"WHERE id = ?", id)

	if err != nil {
		log.Error("[DATABASE]", err)
	}

	return style
}

// LoadEquipment Database Call
func LoadEquipment(db *sqlx.DB, id int32) inventory.Equipment {
	var equip = inventory.Equipment{}
	equip.Init()

	var rows, err = db.Queryx(
		"SELECT kind, serials, opt, slot, expire "+
			"FROM characters_equipment "+
			"WHERE id = ?", id)

	// iterate over each row
	for rows.Next() {
		var i = inventory.Item{}
		var err2 = rows.StructScan(&i)

		if err2 == nil {
			equip.Set(i.Slot, i)
		} else {
			log.Error("[DATABASE]", err2)
			return equip
		}
	}

	if err != nil {
		log.Error("[DATABASE]", err)
	}

	return equip
}

// CreateCharacter RPC Call
func CreateCharacter(_ *rpc.Client, r *character.CreateReq, s *character.CreateRes) error {
	var db = g_DatabaseManager.Get(r.Server)
	var res = character.CreateRes{}

	var c = r.Character
	var cs = c.Style

	// check name
	var name = 0
	db.Get(&name, "SELECT id FROM characters WHERE name = ? LIMIT 1", c.Name)
	if name > 0 {
		res.Result = character.NameInUse
		*s = res
		return nil
	}

	// check battle style
	if cs.BattleStyle > 6 {
		res.Result = character.DBError
		*s = res
		return nil
	}

	// get initial data
	var init = g_DataLoader.BattleStyles[cs.BattleStyle-1]
	var l = init.Location
	var st = init.Stats

	// set data
	c.World = byte(l["world"])
	c.X = byte(l["x"])
	c.Y = byte(l["y"])
	c.CurrentHP = uint16(st["hp"])
	c.MaxHP = uint16(st["hp"])
	c.CurrentMP = uint16(st["mp"])
	c.MaxMP = uint16(st["mp"])
	c.STR = uint32(st["str"])
	c.INT = uint32(st["int"])
	c.DEX = uint32(st["dex"])
	c.Created = time.Now()

	var sql = "INSERT INTO characters ("
	sql += "id, name, world, x, y, gender, hair, color, face, battle_style, current_hp,"
	sql += "max_hp, current_mp, max_mp, str_stat, int_stat, dex_stat, created"
	sql += ") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	db.MustExec(sql, c.Id, c.Name, c.World, c.X, c.Y, cs.Gender, cs.HairStyle,
		cs.HairColor, cs.Face, cs.BattleStyle, c.CurrentHP, c.MaxHP, c.CurrentMP,
		c.MaxMP, c.STR, c.INT, c.DEX, time.Now(),
	)

	// create equipment
	c.Equipment = inventory.Equipment{}
	c.Equipment.Init()

	for key, value := range init.Equipment {
		value.Slot = inventory.MapEquipment(key)
		c.Equipment.Set(value.Slot, value)
		db.MustExec("INSERT INTO characters_equipment "+
			"(id, kind, serials, opt, slot, expire) VALUES (?, ?, ?, ?, ?, ?)",
			c.Id, value.Kind, value.Serials, value.Option, value.Slot, value.Expire,
		)
	}

	// create inventory
	for key, value := range init.Inventory {
		value.Slot = uint16(key)
		db.MustExec("INSERT INTO characters_inventory "+
			"(id, kind, serials, opt, slot, expire) VALUES (?, ?, ?, ?, ?, ?)",
			c.Id, value.Kind, value.Serials, value.Option, value.Slot, value.Expire,
		)
	}

	// create skills
	for key, value := range init.Skills {
		value.Slot = uint16(key)
		db.MustExec("INSERT INTO characters_skills "+
			"(id, skill, level, slot) VALUES (?, ?, ?, ?)",
			c.Id, value.Id, value.Level, value.Slot,
		)
	}

	// create skill links
	for key, value := range init.Links {
		value.Slot = uint16(key)
		db.MustExec("INSERT INTO characters_quickslots "+
			"(id, skill, slot) VALUES (?, ?, ?)",
			c.Id, value.Skill, value.Slot,
		)
	}

	// create lobby metadata if doesn't exist
	var account = r.Id >> 3
	db.MustExec("INSERT IGNORE INTO lobby_metadata (id) VALUE (?)", account)

	res.Result = character.Success
	res.Character = c

	*s = res
	return nil
}

// DeleteCharacter RPC Call
func DeleteCharacter(_ *rpc.Client, r *character.DeleteReq, s *character.DeleteRes) error {
	var db = g_DatabaseManager.Get(r.Server)
	var res = character.DeleteRes{}

	if db == nil {
		*s = res
		return nil
	}

	db.MustExec("DELETE FROM characters_equipment WHERE id = ?", r.CharId)
	db.MustExec("DELETE FROM characters_inventory WHERE id = ?", r.CharId)
	db.MustExec("DELETE FROM characters_quickslots WHERE id = ?", r.CharId)
	db.MustExec("DELETE FROM characters_skills WHERE id = ?", r.CharId)
	db.MustExec("DELETE FROM characters WHERE id = ?", r.CharId)

	res.Result = character.Success

	*s = res
	return nil
}

// SetSlotOrder RPC Call
func SetSlotOrder(_ *rpc.Client, r *character.SetOrderReq, s *character.SetOrderRes) error {
	var db = g_DatabaseManager.Get(r.Server)
	var res = character.SetOrderRes{}

	if db == nil {
		*s = res
		return nil
	}

	db.MustExec(
		"UPDATE lobby_metadata SET slot_order = ? WHERE id = ?", r.Order, r.Account)

	res.Result = true

	*s = res
	return nil
}

// LoadCharacterData RPC Call
func LoadCharacterData(c *rpc.Client, r *character.DataReq, s *character.DataRes) error {
	var db = g_DatabaseManager.Get(r.Server)
	var res = character.DataRes{}

	if db == nil {
		*s = res
		return nil
	}

	// load data
	res.Inventory = LoadInventory(db, r.Id)
	res.Skills = LoadSkills(db, r.Id)
	res.Links = LoadLinks(db, r.Id)

	*s = res
	return nil
}

func MoveItemEquToEqu(c *rpc.Client, r *character.ItemEquipReq, s *character.ItemEquipRes) error {
	var db = g_DatabaseManager.Get(r.Server)

	var inv = inventory.Inventory{}
	inv.Init()

	var equ = inventory.Equipment{}
	equ.Init()

	var item inventory.Item
	err := db.QueryRow("SELECT kind, serials, opt, slot, expire FROM characters_inventory WHERE id = ? AND slot = ?", r.Id, r.DeleteSlot).Scan(&item.Kind, &item.Serials, &item.Option, &item.Slot, &item.Expire)

	//log.Printf("EQUIP ITEM ", item)
	//log.Printf("EQUIP [DATABASE]", err)
	//log.Printf("EQUIP Request ", r)

	if err == nil {
		inv.Remove(item.Slot)
		var _, err3 = db.Queryx(
			"DELETE "+
				"FROM characters_inventory "+
				"WHERE id = ? AND slot = ?", r.Id, item.Slot)
		if err3 != nil {
			log.Printf("[DATABASE]", err3)
		}
		item.Slot = r.CreateSlot
		equ.Set(item.Slot, item)
		db.MustExec("INSERT INTO characters_equipment "+
			"(id, kind, serials, opt, slot, expire) VALUES (?, ?, ?, ?, ?, ?)",
			r.Id, item.Kind, item.Serials, item.Option, item.Slot, item.Expire,
		)

		s.ItemKind = uint32(item.Kind)
	}

	if err != nil {
		log.Printf("[DATABASE]", err)
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
	err := db.QueryRow("SELECT kind, serials, opt, slot, expire FROM characters_equipment WHERE id = ? AND slot = ?", r.Id, r.DeleteSlot).Scan(&item.Kind, &item.Serials, &item.Option, &item.Slot, &item.Expire)

	//log.Printf("UNEQUIP ITEM ", item)
	//log.Printf("UNEQUIP [DATABASE]", err)
	//log.Printf("UNEQUIP Request ", r)

	if err == nil {
		equ.Remove(item.Slot)
		var _, err3 = db.Queryx(
			"DELETE "+
				"FROM characters_equipment "+
				"WHERE id = ? AND slot = ?", r.Id, item.Slot)
		if err3 != nil {
			log.Printf("[DATABASE]", err3)
		}
		item.Slot = r.CreateSlot
		inv.Set(item.Slot, item)
		db.MustExec("INSERT INTO characters_inventory "+
			"(id, kind, serials, opt, slot, expire) VALUES (?, ?, ?, ?, ?, ?)",
			r.Id, item.Kind, item.Serials, item.Option, item.Slot, item.Expire,
		)

		s.ItemKind = uint32(item.Kind)
	}

	if err != nil {
		log.Printf("[DATABASE]", err)
	}

	return nil
}

func MoveItemInvToEqu(c *rpc.Client, r *character.ItemEquipReq, s *character.ItemEquipRes) error {
	var db = g_DatabaseManager.Get(r.Server)

	var item inventory.Item
	err := db.QueryRow("SELECT kind, serials, opt, slot, expire FROM characters_inventory WHERE id = ? AND slot = ?", r.Id, r.DeleteSlot).Scan(&item.Kind, &item.Serials, &item.Option, &item.Slot, &item.Expire)

	//log.Printf("UNEQUIP ITEM ", item)
	//log.Printf("UNEQUIP [DATABASE]", err)
	//log.Printf("UNEQUIP Request ", r)

	if err == nil {
		db.MustExec("UPDATE characters_inventory SET slot = ? WHERE id = ? AND slot = ?", r.CreateSlot, r.Id, item.Slot)
		item.Slot = r.CreateSlot

		s.ItemKind = uint32(item.Kind)
	}

	if err != nil {
		log.Printf("[DATABASE]", err)
	}

	return nil
}

func MoveItemInvToInv(c *rpc.Client, r *character.ItemEquipReq, s *character.ItemEquipRes) error {
	var db = g_DatabaseManager.Get(r.Server)

	var item inventory.Item
	err := db.QueryRow("SELECT kind, serials, opt, slot, expire FROM characters_equipment WHERE id = ? AND slot = ?", r.Id, r.DeleteSlot).Scan(&item.Kind, &item.Serials, &item.Option, &item.Slot, &item.Expire)

	//log.Printf("UNEQUIP ITEM ", item)
	//log.Printf("UNEQUIP [DATABASE]", err)
	//log.Printf("UNEQUIP Request ", r)

	if err == nil {
		db.MustExec("UPDATE characters_equipment SET slot = ? WHERE id = ? AND slot = ?", r.CreateSlot, r.Id, item.Slot)
		item.Slot = r.CreateSlot

		s.ItemKind = uint32(item.Kind)
	}

	if err != nil {
		log.Printf("[DATABASE]", err)
	}

	return nil
}

// LoadInventory Database Call
func LoadInventory(db *sqlx.DB, id int32) inventory.Inventory {
	var inv = inventory.Inventory{}
	inv.Init()

	var rows, err = db.Queryx(
		"SELECT kind, serials, opt, slot, expire "+
			"FROM characters_inventory "+
			"WHERE id = ?", id)

	// iterate over each row
	for rows.Next() {
		var i = inventory.Item{}
		var err2 = rows.StructScan(&i)

		if err2 == nil {
			inv.Set(i.Slot, i)
		} else {
			log.Error("[DATABASE]", err2)
			return inv
		}
	}

	if err != nil {
		log.Error("[DATABASE]", err)
	}

	return inv
}

// LoadSkills Database Call
func LoadSkills(db *sqlx.DB, id int32) skills.SkillList {
	var list = skills.SkillList{}
	list.Init()

	var rows, err = db.Queryx(
		"SELECT skill, level, slot "+
			"FROM characters_skills "+
			"WHERE id = ?", id)

	// iterate over each row
	for rows.Next() {
		var s = skills.Skill{}
		var err2 = rows.StructScan(&s)

		if err2 == nil {
			list.Set(s.Slot, s)
		} else {
			log.Error("[DATABASE]", err2)
			return list
		}
	}

	if err != nil {
		log.Error("[DATABASE]", err)
	}

	return list
}

// LoadLinks Database Call
func LoadLinks(db *sqlx.DB, id int32) skills.Links {
	var list = skills.Links{}
	list.Init()

	var rows, err = db.Queryx(
		"SELECT skill, slot "+
			"FROM characters_quickslots "+
			"WHERE id = ?", id)

	// iterate over each row
	for rows.Next() {
		var l = skills.Link{}
		var err2 = rows.StructScan(&l)

		if err2 == nil {
			list.Set(l.Slot, l)
		} else {
			log.Error("[DATABASE]", err2)
			return list
		}
	}

	if err != nil {
		log.Error("[DATABASE]", err)
	}

	return list
}
