package rpc

import (
	"context"

	"github.com/ubis/Freya/share/models/character"
	"github.com/ubis/Freya/share/models/inventory"
	"github.com/ubis/Freya/share/rpc"
)

func UnEquipItem(c *rpc.Client, r *character.ItemMoveReq, s *character.ItemMoveRes) error {
	var db = g_DatabaseManager.Get(r.Server)

	var item inventory.Item
	query := db.QueryRow("SELECT kind, serials, opt, slot, expire "+
		"FROM characters_equipment WHERE id = ? AND slot = ?", r.Id, r.DeleteSlot)
	err := query.Scan(&item.Kind, &item.Serials, &item.Option, &item.Slot, &item.Expire)
	if err != nil {
		s.Result = err
		return err
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		s.Result = err
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM characters_equipment WHERE id = ? AND slot = ?", r.Id, item.Slot)
	if err != nil {
		_ = tx.Rollback()
		s.Result = err
		return err
	}

	item.Slot = r.CreateSlot

	_, err = tx.ExecContext(ctx, "INSERT INTO characters_inventory "+
		"(id, kind, serials, opt, slot, expire) "+
		"VALUES (?, ?, ?, ?, ?, ?)",
		r.Id, item.Kind, item.Serials, item.Option, r.CreateSlot, item.Expire)
	if err != nil {
		_ = tx.Rollback()
		s.Result = err
		return err
	}

	err = tx.Commit()
	if err != nil {
		s.Result = err
		return err
	}

	s.Item = item
	s.Result = nil

	return nil
}

func EquipItem(c *rpc.Client, r *character.ItemMoveReq, s *character.ItemMoveRes) error {
	var db = g_DatabaseManager.Get(r.Server)

	var item inventory.Item
	query := db.QueryRow("SELECT kind, serials, opt, slot, expire "+
		"FROM characters_inventory WHERE id = ? AND slot = ?", r.Id, r.DeleteSlot)
	err := query.Scan(&item.Kind, &item.Serials, &item.Option, &item.Slot, &item.Expire)
	if err != nil {
		s.Result = err
		return err
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		s.Result = err
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM characters_inventory WHERE id = ? AND slot = ?", r.Id, item.Slot)
	if err != nil {
		_ = tx.Rollback()
		s.Result = err
		return err
	}

	item.Slot = r.CreateSlot

	_, err = tx.ExecContext(ctx, "INSERT INTO characters_equipment "+
		"(id, kind, serials, opt, slot, expire) "+
		"VALUES (?, ?, ?, ?, ?, ?)",
		r.Id, item.Kind, item.Serials, item.Option, r.CreateSlot, item.Expire)
	if err != nil {
		_ = tx.Rollback()
		s.Result = err
		return err
	}

	err = tx.Commit()
	if err != nil {
		s.Result = err
		return err
	}

	s.Item = item
	s.Result = nil

	return nil
}

func ChangeEquipItemSlot(c *rpc.Client, r *character.ItemMoveReq, s *character.ItemMoveRes) error {
	var db = g_DatabaseManager.Get(r.Server)

	var item inventory.Item
	query := db.QueryRow("SELECT kind, serials, opt, slot, expire "+
		"FROM characters_equipment WHERE id = ? AND slot = ?", r.Id, r.DeleteSlot)
	err := query.Scan(&item.Kind, &item.Serials, &item.Option, &item.Slot, &item.Expire)
	if err != nil {
		s.Result = err
		return err
	}

	affected, err := db.MustExec("UPDATE characters_equipment SET slot = ? "+
		" WHERE id = ? AND slot = ?", r.CreateSlot, r.Id, item.Slot).RowsAffected()
	if err == nil && affected > 0 {
		s.Item = item
		s.Result = nil

		return nil
	}

	s.Result = err
	return nil
}

func ChangeInventoryItemSlot(c *rpc.Client, r *character.ItemMoveReq, s *character.ItemMoveRes) error {
	var db = g_DatabaseManager.Get(r.Server)

	var item inventory.Item
	query := db.QueryRow("SELECT kind, serials, opt, slot, expire "+
		"FROM characters_inventory WHERE id = ? AND slot = ?", r.Id, r.DeleteSlot)
	err := query.Scan(&item.Kind, &item.Serials, &item.Option, &item.Slot, &item.Expire)
	if err != nil {
		s.Result = err
		return err
	}

	affected, err := db.MustExec("UPDATE characters_inventory SET slot = ? "+
		"WHERE id = ? AND slot = ?", r.CreateSlot, r.Id, item.Slot).RowsAffected()
	if err == nil && affected > 0 {
		s.Item = item
		s.Result = nil

		return nil
	}

	s.Result = err
	return nil
}

func PickItem(c *rpc.Client, r *character.ItemPickRequest, s *character.ItemPickResponse) error {
	var db = g_DatabaseManager.Get(r.Server)

	s.Result = false

	_, err := db.MustExec("INSERT INTO characters_inventory "+
		"(id, kind, serials, opt, slot, expire) "+
		"VALUES (?, ?, ?, ?, ?, ?)",
		r.Id, r.Item.Kind, r.Item.Serials, r.Item.Option, r.Item.Slot, r.Item.Expire).RowsAffected()
	if err != nil {
		return err
	}

	s.Result = true

	return nil
}

func DropItem(c *rpc.Client, r *character.ItemPickRequest, s *character.ItemPickResponse) error {
	var db = g_DatabaseManager.Get(r.Server)

	s.Result = false

	_, err := db.MustExec(
		"DELETE FROM characters_inventory WHERE id = ? AND slot = ?",
		r.Id, r.Item.Slot).RowsAffected()
	if err != nil {
		return err
	}

	s.Result = true

	return nil
}

func SwapItem(c *rpc.Client, r *character.ItemSwapRequest, s *character.ItemPickResponse) error {
	var db = g_DatabaseManager.Get(r.Server)

	s.Result = false

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// the rollback will be ignored if the tx has been committed later in the function
	defer tx.Rollback()

	// slots were swapped before
	oldSlot := r.New.Slot
	newSlot := r.Old.Slot

	// since id + slot is primary key, we need to switch to temp slot
	// slot is uint16, so use it's max value
	tempSlot := 65535

	_, err = tx.Exec(
		"UPDATE characters_inventory SET slot = ? WHERE id = ? AND slot = ?",
		tempSlot, r.Id, oldSlot,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"UPDATE characters_inventory SET slot = ? WHERE id = ? AND slot = ?",
		oldSlot, r.Id, newSlot,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"UPDATE characters_inventory SET slot = ? WHERE id = ? AND slot = ?",
		newSlot, r.Id, tempSlot,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	s.Result = true

	return nil
}
