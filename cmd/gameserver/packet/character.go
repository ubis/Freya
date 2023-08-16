package packet

import (
	"bytes"

	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/character"
	"github.com/ubis/Freya/share/models/subpasswd"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

// NewTargetUser Packet
func NewTargetUser(session *network.Session, reader *network.Reader) {
	sessionId := reader.ReadUint16()

	pSession := g_NetworkManager.GetSession(sessionId)
	ctx, err := context.Parse(pSession)
	if err != nil {
		log.Error(err.Error())
		return
	}

	ctx.Mutex.Lock()
	currentHP, maxHP := ctx.Char.CurrentHP, ctx.Char.MaxHP
	ctx.Mutex.RUnlock()

	var packet = network.NewWriter(NEW_TARGET_USER)

	packet.WriteByte(0x00)
	packet.WriteInt16(currentHP)
	packet.WriteInt16(maxHP)

	session.Send(packet)
}

// GetMyChartr Packet
func GetMyChartr(session *network.Session, reader *network.Reader) {
	if !session.Data.Verified {
		log.Error("Unauthorized connection from", session.GetEndPnt())
		session.Close()
		return
	}

	// fetch subpassword
	var req = subpasswd.FetchReq{session.Data.AccountId}
	var res = subpasswd.FetchRes{}
	g_RPCHandler.Call(rpc.FetchSubPassword, req, &res)

	session.Data.SubPassword = &res.Details

	var subpasswdExist = 0
	if res.Password != "" {
		subpasswdExist = 1
	}

	// fetch characters
	var reqList = character.ListReq{session.Data.AccountId, byte(g_ServerSettings.ServerId)}
	var resList = character.ListRes{}
	g_RPCHandler.Call(rpc.LoadCharacters, reqList, &resList)

	session.Data.CharacterList = resList.List

	var packet = network.NewWriter(GETMYCHARTR)
	packet.WriteInt32(subpasswdExist)
	packet.WriteBytes(make([]byte, 10))
	packet.WriteInt32(resList.LastId)
	packet.WriteInt32(resList.SlotOrder)

	for i := 0; i < len(resList.List); i++ {
		var char = resList.List[i]
		packet.WriteInt32(char.Id)
		packet.WriteInt64(char.Created.Unix())
		packet.WriteUint32(char.Style.Get())
		packet.WriteUint32(char.Level)
		packet.WriteByte(char.SwordRank)
		packet.WriteByte(char.MagicRank)
		packet.WriteInt16(0x00) // padding for skill ranks
		packet.WriteUint64(char.Alz)
		packet.WriteByte(char.Nation)
		packet.WriteByte(char.World)
		packet.WriteUint16(char.X)
		packet.WriteUint16(char.Y)
		packet.WriteBytes(char.Equipment.SerializeKind())
		packet.WriteBytes(make([]byte, 88))
		packet.WriteByte(len(char.Name) + 1)
		packet.WriteString(char.Name + "\x00")
	}

	session.Send(packet)
}

// NewMyChartr Packet
func NewMyChartr(session *network.Session, reader *network.Reader) {
	var style = reader.ReadUint32()
	var _ = reader.ReadByte() // beginner join guild
	var slot = reader.ReadByte()
	var nameLength = reader.ReadByte()
	var name = string(bytes.Trim(reader.ReadBytes(int(nameLength)), "\x00"))

	var charId = session.Data.AccountId*8 + int32(slot)
	var newStyle = character.Style{}
	newStyle.Set(style)

	var packet = network.NewWriter(NEWMYCHARTR)

	if !newStyle.Verify() || slot > 5 || nameLength > 16 {
		// invalid style, slot or nameLength
		packet.WriteInt32(0x00)
		packet.WriteByte(character.NowAllowed)

		session.Send(packet)
		return
	}

	// check if slot is used
	var charList = session.Data.CharacterList
	for i := 0; i < len(charList); i++ {
		if charList[i].Id == charId {
			packet.WriteInt32(0x00)
			packet.WriteByte(character.SlotInUse)

			session.Send(packet)
			return
		}
	}

	var req = character.CreateReq{
		byte(g_ServerSettings.ServerId),
		character.Character{Id: charId, Name: name, Style: newStyle},
	}
	var res = character.CreateRes{}
	g_RPCHandler.Call(rpc.CreateCharacter, req, &res)

	if res.Result == character.Success {
		packet.WriteInt32(charId)
		packet.WriteByte(res.Result + newStyle.BattleStyle)
		// update character with it's data
		session.Data.CharacterList = append(session.Data.CharacterList, res.Character)
	} else {
		packet.WriteInt32(0x00)
		packet.WriteByte(res.Result)
	}

	session.Send(packet)
}

// DelMyChartr Packet
func DelMyChartr(session *network.Session, reader *network.Reader) {
	var charId = reader.ReadInt32()

	// if password wasn't verified
	if !session.Data.CharVerified {
		return
	}

	// if subpasswd wasn't verified
	if len(session.Data.SubPassword.Password) > 0 && !session.Data.SubPassword.Verified {
		return
	}

	// verify character id
	if charId>>3 != session.Data.AccountId {
		return
	}

	var req = character.DeleteReq{byte(g_ServerSettings.ServerId), charId}
	var res = character.DeleteRes{}
	g_RPCHandler.Call(rpc.DeleteCharacter, req, &res)

	if res.Result == character.Success {
		// reset character delete passwd verification
		session.Data.CharVerified = false

		// reset character delete subpasswd verification
		session.Data.SubPassword.Verified = false

		var l = &session.Data.CharacterList

		// remove character from the list
		for key, value := range *l {
			if value.Id == charId {
				*l = append((*l)[:key], (*l)[key+1:]...)
				break
			}
		}
	}

	var packet = network.NewWriter(DELMYCHARTR)
	packet.WriteByte(res.Result + 1)
	packet.WriteByte(0x00)

	session.Send(packet)
}

// SetCharacterSlotOrder Packet
func SetCharacterSlotOrder(session *network.Session, reader *network.Reader) {
	var order = reader.ReadInt32()

	var req = character.SetOrderReq{
		byte(g_ServerSettings.ServerId),
		session.Data.AccountId,
		order,
	}
	var res = character.SetOrderRes{}
	g_RPCHandler.Call(rpc.SetSlotOrder, req, &res)

	var packet = network.NewWriter(SET_CHAR_SLOT_ORDER)
	packet.WriteByte(0x01)

	session.Send(packet)
}

func notifyChangeStyle(session *network.Session) {
	ctx, err := context.Parse(session)
	if err != nil {
		log.Error(err.Error())
		return
	}

	ctx.Mutex.RLock()
	id := ctx.Char.Id
	style := ctx.Char.Style
	liveStyle := ctx.Char.LiveStyle
	ctx.Mutex.RUnlock()

	pkt := network.NewWriter(NFY_CHANGESTYLE)
	pkt.WriteInt32(id)
	pkt.WriteInt32(style.Get())
	pkt.WriteInt32(liveStyle)
	pkt.WriteInt32(0)
	pkt.WriteInt16(0)

	ctx.World.BroadcastSessionPacket(session, pkt)
}

func ChangeStyle(session *network.Session, reader *network.Reader) {
	_ = reader.ReadInt32() // style
	liveStyle := reader.ReadInt32()
	_ = reader.ReadInt32() // buffFlag?
	_ = reader.ReadInt16() // actionFlag?

	ctx, err := context.Parse(session)
	if err != nil {
		log.Error(err.Error())
		return
	}

	ctx.Mutex.Lock()
	ctx.Char.LiveStyle = liveStyle
	ctx.Mutex.Unlock()

	pkt := network.NewWriter(CHANGESTYLE)
	pkt.WriteByte(1)

	session.Send(pkt)

	notifyChangeStyle(session)
}

func SkillToActs(session *network.Session, reader *network.Reader) {
	target := reader.ReadInt32() // self char id
	action := reader.ReadUint16()
	x := reader.ReadByte()
	y := reader.ReadByte()

	id, err := context.GetCharId(session)
	if err != nil {
		log.Error(err.Error())
		return
	}

	ctx, err := context.Parse(session)
	if err != nil {
		log.Error(err.Error())
		return
	}

	pkt := network.NewWriter(NFY_SKILLTOACTS)
	pkt.WriteInt32(id)
	pkt.WriteInt32(target)
	pkt.WriteUint16(action)
	pkt.WriteByte(x)
	pkt.WriteByte(y)

	ctx.World.BroadcastSessionPacket(session, pkt)
}
