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
