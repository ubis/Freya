package rpc

import (
	"errors"

	"github.com/ubis/Freya/share/models/skills"
	"github.com/ubis/Freya/share/rpc"
)

func QuickLinkSet(c *rpc.Client, r *skills.QuickLinkRequest, s *skills.QuickLinkResponse) error {
	var db = g_DatabaseManager.Find(c)

	s.Result = false

	_, err := db.MustExec(
		"INSERT INTO characters_quickslots "+
			"(id, skill, slot) "+
			"VALUES (?, ?, ?)",
		r.Id, r.NewLink.Skill, r.NewLink.Slot).RowsAffected()
	if err != nil {
		return err
	}

	s.Result = true

	return nil
}

func QuickLinkRemove(c *rpc.Client, r *skills.QuickLinkRequest, s *skills.QuickLinkResponse) error {
	var db = g_DatabaseManager.Find(c)

	s.Result = false

	_, err := db.MustExec(
		"DELETE FROM characters_quickslots WHERE id = ? AND slot = ?",
		r.Id, r.OldLink.Slot).RowsAffected()
	if err != nil {
		return err
	}

	s.Result = true

	return nil
}

func QuickLinkSwap(c *rpc.Client, r *skills.QuickLinkRequest, s *skills.QuickLinkResponse) error {
	var db = g_DatabaseManager.Find(c)

	s.Result = false

	if r.OldLink == nil {
		return errors.New("source link is not set")
	}

	if r.NewLink == nil {
		return errors.New("target link is not set")
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// the rollback will be ignored if the tx has been committed later in the function
	defer tx.Rollback()

	// slots were swapped before
	newSlot := r.OldLink.Slot
	oldSlot := r.NewLink.Slot

	// since id + slot is primary key, we need to switch to temp slot
	// slot is uint16, so use it's max value
	tempSlot := 65535

	_, err = tx.Exec(
		"UPDATE characters_quickslots SET slot = ? WHERE id = ? AND slot = ?",
		tempSlot, r.Id, oldSlot,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"UPDATE characters_quickslots SET slot = ? WHERE id = ? AND slot = ?",
		oldSlot, r.Id, newSlot,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"UPDATE characters_quickslots SET slot = ? WHERE id = ? AND slot = ?",
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
