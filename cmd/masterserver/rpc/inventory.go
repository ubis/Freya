package rpc

import (
	"errors"

	"github.com/ubis/Freya/share/models/inventory"
	"github.com/ubis/Freya/share/rpc"
)

func EquipItem(c *rpc.Client, r *inventory.ItemRequest, s *inventory.ItemResponse) error {
	var db = g_DatabaseManager.Get(r.Server)

	s.Result = false

	_, err := db.MustExec("INSERT INTO characters_equipment "+
		"(id, kind, serials, opt, slot, expire) "+
		"VALUES (?, ?, ?, ?, ?, ?)",
		r.Id, r.Item.Kind, r.Item.Serials, r.Item.Option, r.Item.Slot, r.Item.Expire).RowsAffected()
	if err != nil {
		return err
	}

	s.Result = true

	return nil
}

func UnEquipItem(c *rpc.Client, r *inventory.ItemRequest, s *inventory.ItemResponse) error {
	var db = g_DatabaseManager.Get(r.Server)

	s.Result = false

	_, err := db.MustExec(
		"DELETE FROM characters_equipment WHERE id = ? AND slot = ?",
		r.Id, r.Item.Slot).RowsAffected()
	if err != nil {
		return err
	}

	s.Result = true

	return nil
}

func MoveEquipmentItem(c *rpc.Client, r *inventory.ItemRequest, s *inventory.ItemResponse) error {
	var db = g_DatabaseManager.Get(r.Server)

	s.Result = false

	_, err := db.MustExec(
		"UPDATE characters_equipment SET slot = ? "+
			"WHERE id = ? AND slot = ?",
		r.NewItem.Slot, r.Id, r.Item.Slot).RowsAffected()
	if err != nil {
		return err
	}

	s.Result = true

	return nil
}

func AddItem(c *rpc.Client, r *inventory.ItemRequest, s *inventory.ItemResponse) error {
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

func RemoveItem(c *rpc.Client, r *inventory.ItemRequest, s *inventory.ItemResponse) error {
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

func SwapItem(c *rpc.Client, r *inventory.ItemRequest, s *inventory.ItemResponse) error {
	var db = g_DatabaseManager.Get(r.Server)

	s.Result = false

	if r.NewItem == nil {
		return errors.New("target item is not set!")
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// the rollback will be ignored if the tx has been committed later in the function
	defer tx.Rollback()

	// slots were swapped before
	newSlot := r.Item.Slot
	oldSlot := r.NewItem.Slot

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

func MoveItem(c *rpc.Client, r *inventory.ItemRequest, s *inventory.ItemResponse) error {
	var db = g_DatabaseManager.Get(r.Server)

	s.Result = false

	_, err := db.MustExec(
		"UPDATE characters_inventory SET slot = ? "+
			"WHERE id = ? AND slot = ?",
		r.NewItem.Slot, r.Id, r.Item.Slot).RowsAffected()
	if err != nil {
		return err
	}

	s.Result = true

	return nil
}
