package packet

import (
	"bytes"

	"github.com/ubis/Freya/share/models/character"
	"github.com/ubis/Freya/share/models/subpasswd"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

// NewTargetUser Packet
func NewTargetUser(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	userIdx := reader.ReadUint16()

	player := session.FindPlayerByIndex(userIdx)
	if player == nil {
		session.LogErrorf("Invalid session index requested: %d", userIdx)
		return
	}

	currentHP, maxHP := player.Character.GetHealth()

	pkt := network.NewWriter(CSCNewTargetUser)
	pkt.WriteByte(0x00)
	pkt.WriteInt16(currentHP)
	pkt.WriteInt16(maxHP)

	session.Send(pkt)
}

// GetMyChar Packet
func GetMyChar(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	// fetch subpassword
	req := subpasswd.FetchReq{Account: session.Account}
	res := subpasswd.FetchRes{}
	session.RPC.Call(rpc.FetchSubPassword, &req, &res)

	session.SubPassword = &res.Details

	subpasswdExist := 0
	if session.ServerConfig.IgnoreSubPassword || res.Password != "" {
		subpasswdExist = 1
	}

	// fetch characters
	reqList := character.ListReq{Account: session.Account}
	resList := character.ListRes{}
	session.RPC.Call(rpc.LoadCharacters, &reqList, &resList)

	session.Characters = resList.List

	pkt := network.NewWriter(CSCGetMyChar)
	pkt.WriteInt32(subpasswdExist)
	pkt.WriteBytes(make([]byte, 10))
	pkt.WriteInt32(resList.LastId)
	pkt.WriteInt32(resList.SlotOrder)

	for _, v := range resList.List {
		pkt.WriteInt32(v.Id)
		pkt.WriteInt64(v.Created.Unix())
		pkt.WriteUint32(v.Style.Get())
		pkt.WriteUint32(v.Level)
		pkt.WriteByte(v.SwordRank)
		pkt.WriteByte(v.MagicRank)
		pkt.WriteInt16(0x00) // padding for skill ranks
		pkt.WriteUint64(v.Alz)
		pkt.WriteByte(v.Nation)
		pkt.WriteByte(v.World)
		pkt.WriteUint16(v.X)
		pkt.WriteUint16(v.Y)
		pkt.WriteBytes(v.Equipment.SerializeKind())
		pkt.WriteBytes(make([]byte, 88))
		pkt.WriteByte(len(v.Name) + 1)
		pkt.WriteString(v.Name + "\x00")
	}

	session.Send(pkt)
}

// NewMyChar Packet
func NewMyChar(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	style := reader.ReadUint32()
	_ = reader.ReadByte() // beginner join guild
	slot := int32(reader.ReadByte())
	nameLength := reader.ReadByte()
	name := string(bytes.Trim(reader.ReadBytes(int(nameLength)), "\x00"))

	charId := session.Account*8 + slot
	newStyle := character.Style{}
	newStyle.Set(style)

	pkt := network.NewWriter(CSCNewMyChar)

	if !newStyle.Verify() ||
		slot > MaxCharacterSlot ||
		nameLength > MaxCharacterNameLength {

		session.LogErrorf("Invalid character create data: "+
			"style: %d ; "+
			"slot: %d of %d ; "+
			"name length: %d of %d",
			style, slot, MaxCharacterSlot, nameLength, MaxCharacterNameLength)

		// invalid style, slot or name length
		pkt.WriteInt32(0x00)
		pkt.WriteByte(character.NotAllowed)

		session.Send(pkt)
		return
	}

	// check if slot is used
	for _, v := range session.Characters {
		if v.Id != charId {
			continue
		}

		session.LogErrorf("Player is trying to create a character "+
			" on already occupied slot: %d ", slot)

		pkt.WriteInt32(0x00)
		pkt.WriteByte(character.SlotInUse)
		session.Send(pkt)
		return
	}

	req := character.CreateReq{
		Character: character.Character{
			Id:    charId,
			Name:  name,
			Style: newStyle,
		},
	}
	res := character.CreateRes{}
	session.RPC.Call(rpc.CreateCharacter, &req, &res)

	if res.Result == character.Success {
		pkt.WriteInt32(charId)
		pkt.WriteByte(res.Result + newStyle.BattleStyle)

		// update character with it's data
		session.Characters = append(session.Characters, &res.Character)
	} else {
		pkt.WriteInt32(0x00)
		pkt.WriteByte(res.Result)
	}

	session.Send(pkt)
}

// DelMyChar Packet
func DelMyChar(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	charId := reader.ReadInt32()

	// if password wasn't verified
	if !session.PasswordVerified {
		session.LogErrorf("Player is trying to delete a character: %d "+
			"without password verification", charId)
		return
	}

	// if subpasswd wasn't verified
	sub := session.SubPassword
	if len(sub.Password) > 0 && !sub.Verified {
		session.LogErrorf("Player is trying to delete a character: %d "+
			"without subpassword verification", charId)
		return
	}

	// verify character id
	if charId>>3 != session.Account {
		session.LogErrorf("Player is trying to delete an invalid character: %d",
			charId)
		return
	}

	req := character.DeleteReq{
		CharId: charId,
	}
	res := character.DeleteRes{}
	session.RPC.Call(rpc.DeleteCharacter, &req, &res)

	if res.Result == character.Success {
		// reset character delete passwd verification
		session.PasswordVerified = false

		// reset character delete subpasswd verification
		sub.Verified = false

		l := &session.Characters

		// remove character from the list
		for key, value := range *l {
			if value.Id == charId {
				*l = append((*l)[:key], (*l)[key+1:]...)
				break
			}
		}
	}

	pkt := network.NewWriter(CSCDelMyChar)
	pkt.WriteByte(res.Result)
	pkt.WriteByte(0x00)

	session.Send(pkt)
}

// SetCharacterSlotOrder Packet
func SetCharacterSlotOrder(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	order := reader.ReadInt32()

	req := character.SetOrderReq{
		Account: session.Account,
		Order:   order,
	}
	res := character.SetOrderRes{}
	session.RPC.Call(rpc.SetSlotOrder, &req, &res)

	pkt := network.NewWriter(CSCSetCharacterSlotOrder)
	pkt.WriteBool(res.Result)

	session.Send(pkt)
}

// ChangeStyle Packet
func ChangeStyle(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	_ = reader.ReadInt32() // style
	liveStyle := reader.ReadInt32()
	_ = reader.ReadInt32() // buffFlag?
	_ = reader.ReadInt16() // actionFlag?

	// todo: verify live style
	session.Character.SetLiveStyle(liveStyle)

	pkt := network.NewWriter(CSCChangeStyle)
	pkt.WriteBool(true)

	session.Send(pkt)

	// notify surrounding players
	style, _ := session.Character.GetStyle()

	pkt = network.NewWriter(NFYChangeStyle)
	pkt.WriteInt32(session.Character.Id)
	pkt.WriteInt32(style.Get())
	pkt.WriteInt32(liveStyle)
	pkt.WriteInt32(0)
	pkt.WriteInt16(0)

	session.World.BroadcastSessionPacket(session.SessionHandler, pkt)
}

// SkillToAction Packet
func SkillToAction(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	target := reader.ReadInt32() // self char id
	action := reader.ReadUint16()
	x := reader.ReadByte()
	y := reader.ReadByte()

	pkt := network.NewWriter(NFYSkillToAction)
	pkt.WriteInt32(session.Character.Id)
	pkt.WriteInt32(target)
	pkt.WriteUint16(action)
	pkt.WriteByte(x)
	pkt.WriteByte(y)

	session.World.BroadcastSessionPacket(session.SessionHandler, pkt)
}

// SkillToUser Packet
func SkillToUser(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	// seems like there are 2 types of messages
	switch reader.Size {
	case 19:
		// astral & style related
		handleStyleSkill(session, reader)
	case 17:
		// dash/fade & movement related
		handleMoveSkill(session, reader)
	}
}

func handleStyleSkill(session *Session, reader *network.Reader) {
	skill := reader.ReadUint16()
	_ = reader.ReadByte() // slot
	unk1 := reader.ReadInt16()
	unk2 := reader.ReadInt32()

	mp, _ := session.Character.GetMana()

	pkt := network.NewWriter(CSCSkillToUser)
	pkt.WriteUint16(skill)
	pkt.WriteUint16(mp)
	pkt.WriteInt16(unk1)
	pkt.WriteInt32(unk2)

	session.Send(pkt)

	// notify surrounding players
	style, liveStyle := session.Character.GetStyle()

	pkt = network.NewWriter(NFYSkillToUser)
	pkt.WriteUint16(skill)
	pkt.WriteUint32(session.Character.Id)
	pkt.WriteUint32(style)
	pkt.WriteByte(liveStyle)
	pkt.WriteByte(0x02)
	pkt.WriteInt16(unk1)
	pkt.WriteInt32(unk2)

	session.World.BroadcastSessionPacket(session.SessionHandler, pkt)
}

func handleMoveSkill(session *Session, reader *network.Reader) {
	skill := reader.ReadUint16()
	_ = reader.ReadByte() // slot
	x := byte(reader.ReadInt16())
	y := byte(reader.ReadInt16())

	mp, _ := session.Character.GetMana()

	pkt := network.NewWriter(CSCSkillToUser)
	pkt.WriteUint16(skill)
	pkt.WriteInt32(0)
	pkt.WriteUint16(mp)

	session.Send(pkt)

	// notify surrounding players
	session.Character.SetPosition(x, y)
	session.Character.SetMovement(x, y, x, y)

	pkt = network.NewWriter(NFYSkillToUser)
	pkt.WriteUint16(skill)
	pkt.WriteUint32(session.Character.Id)
	pkt.WriteUint16(session.GetUserIdx())
	pkt.WriteInt16(0x1000)
	pkt.WriteInt16(x)
	pkt.WriteInt16(y)

	session.World.BroadcastSessionPacket(session.SessionHandler, pkt)
	session.World.AdjustCell(session.SessionHandler)
}

func GetPlayerLevel(session *Session) int {
	// fixme
	return int(session.Character.GetLevel())
}

func SetPlayerLevel(session *Session, level int) {
	// fixme
	session.Character.SetLevel(uint16(level))

	pkt := network.NewWriter(288)
	pkt.WriteByte(1) // 1 = level up; 2 = rank up
	pkt.WriteInt32(session.Character.Id)

	session.World.BroadcastSessionPacket(session.SessionHandler, pkt)

	pkt = network.NewWriter(287)
	pkt.WriteByte(10) // 10 = level up
	for i := 0; i < 14; i++ {
		pkt.WriteByte(0)
	}
	pkt.WriteInt64(level)

	session.Send(pkt)
}
