package character

type Style struct {
    BattleStyle  byte `db:"battle_style"`
    MasteryLevel byte `db:"rank"`
    Face         byte
    HairColor    byte `db:"color"`
    HairStyle    byte `db:"hair"`
    Aura         byte
    Gender       bool
    Helmet       bool `db:"show_helmet"`
}

// Deserializes style from uint32 to struct
func (s *Style) Set(style uint32) {
    s.BattleStyle  = byte(style & 0x07)
    s.MasteryLevel = byte((style >> 3) & 0x1F)
    s.Face         = byte((style >> 8) & 0x1F)
    s.HairColor    = byte((style >> 13) & 0x0F)
    s.HairStyle    = byte((style >> 17) & 0x1F)
    s.Aura         = byte((style >> 22) & 0x0F)

    if s.Gender = false; ((style >> 26) & 0x01) == 0x01 {
        s.Gender = true
    }

    if s.Helmet = false; ((style >> 27) & 0x01) == 0x01 {
        s.Helmet = true
    }
}

// Serializes style from uint32 to struct
func (s *Style) Get() uint32 {
    var style = uint32(s.BattleStyle)
    style += uint32(s.MasteryLevel) << 3
    style += uint32(s.Face) << 8
    style += uint32(s.HairColor) << 13
    style += uint32(s.HairStyle) << 17
    style += uint32(s.Aura) << 22

    if !s.Gender {
        style += uint32(0) << 26
    } else {
        style += uint32(1) << 26
    }

    if !s.Helmet {
        style += uint32(0) << 27
    } else {
        style += uint32(1) << 27
    }

    return style
}

// Verifies style and if one of parameters are invalid, false will be returned
func (s *Style) Verify() bool {
    s.MasteryLevel = 1
    s.Aura         = 0
    s.Helmet       = false

    // check battle style
    if s.BattleStyle < 1 || s.BattleStyle > 6 {
        return false
    }

    // check face, color and hair
    if s.Face > 3 || s.HairColor > 7 || s.HairStyle > 6  {
        return false
    }

    return true
}