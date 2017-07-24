package rpc

import (
    "time"
    "share/rpc"
    "share/models/character"
    "share/models/inventory"
    "github.com/jmoiron/sqlx"
)

// LoadCharacters RPC Call
func LoadCharacters(_ *rpc.Client, r *character.ListReq, s *character.ListRes) error {
    var db  = g_DatabaseManager.Get(r.Server)
    var res = character.ListRes{make([]character.Character, 0, 6)}

    if db == nil {
        *s = res
        return nil
    }

    var rows, err = db.Queryx(
        "SELECT " +
            "id, name, level, world, x, y, alz, nation, sword_rank, magic_rank, created " +
        "FROM characters " +
        "WHERE id >= ? AND id <= ?", r.Account * 8, r.Account * 8 + 5)

    if err != nil {
        log.Error("[DATABASE]", err)
        *s = res
        return nil
    }

    // iterate over each row
    for rows.Next() {
        var c   = character.Character{}
        var err = rows.StructScan(&c)

        if err == nil {
            c.Style     = LoadStyle(db, c.Id)
            c.Equipment = LoadEquipment(db, c.Id)

            res.List = append(res.List, c)
        } else {
            log.Error("[DATABASE]", err)
            *s = res
            return nil
        }
    }

    *s = res
    return nil
}

// LoadStyle Database Call
func LoadStyle(db *sqlx.DB, id int32) character.Style {
    var style = character.Style{}
    var err   = db.Get(&style,
        "SELECT battle_style, rank, face, color, hair, aura, gender, show_helmet " +
        "FROM characters " +
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
        "SELECT kind, serials, opt, slot, expire " +
        "FROM characters_equipment " +
        "WHERE id = ?", id)

    // iterate over each row
    for rows.Next() {
        var i    = inventory.Item{}
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
    var db  = g_DatabaseManager.Get(r.Server)
    var res = character.CreateRes{}

    var c  = r.Character
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
    var init = g_DataLoader.BattleStyles[cs.BattleStyle - 1]
    var l    = init.Location
    var st   = init.Stats

    // set data
    c.World     = byte(l["world"])
    c.X         = byte(l["x"])
    c.Y         = byte(l["y"])
    c.CurrentHP = uint16(st["hp"])
    c.MaxHP     = uint16(st["hp"])
    c.CurrentMP = uint16(st["mp"])
    c.MaxMP     = uint16(st["mp"])
    c.STR       = uint32(st["str"])
    c.INT       = uint32(st["int"])
    c.DEX       = uint32(st["dex"])
    c.Created   = time.Now()

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
        db.MustExec("INSERT INTO characters_equipment " +
                "(id, kind, serials, opt, slot, expire) VALUES (?, ?, ?, ?, ?, ?)",
            c.Id, value.Kind, value.Serials, value.Option, value.Slot, value.Expire,
        )
    }

    // create inventory
    for key, value := range init.Inventory {
        value.Slot = uint16(key)
        db.MustExec("INSERT INTO characters_inventory " +
            "(id, kind, serials, opt, slot, expire) VALUES (?, ?, ?, ?, ?, ?)",
            c.Id, value.Kind, value.Serials, value.Option, value.Slot, value.Expire,
        )
    }

    // create skills
    for key, value := range init.Skills {
        value.Slot = uint16(key)
        db.MustExec("INSERT INTO characters_skills " +
            "(id, skill, level, slot) VALUES (?, ?, ?, ?)",
            c.Id, value.Id, value.Level, value.Slot,
        )
    }

    // create skill links
    for key, value := range init.Links {
        value.Slot = uint16(key)
        db.MustExec("INSERT INTO characters_quickslots " +
            "(id, skill, slot) VALUES (?, ?, ?)",
            c.Id, value.Slot, value.Slot,
        )
    }

    res.Result    = character.Success
    res.Character = c

    *s = res
    return nil
}