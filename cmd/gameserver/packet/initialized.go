package packet

import (
	"strings"
	"time"

	"github.com/ubis/Freya/share/event"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/character"
	"github.com/ubis/Freya/share/models/server"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
	"github.com/ubis/Freya/share/script"
)

// Initialize Packet
func Initialize(session *Session, reader *network.Reader) {
	if !verifyState(session, StateVerified, reader.Type) {
		return
	}

	charId := reader.ReadInt32()

	// verify char id
	if (charId >> 3) != session.Account {
		session.LogFatalf("Player is trying to use an invalid character: %d",
			charId)
		return
	}

	// fetch characters
	if len(session.Characters) == 0 {
		reqList := character.ListReq{Account: session.Account}
		resList := character.ListRes{}
		session.RPC.Call(rpc.LoadCharacters, &reqList, &resList)

		session.Characters = resList.List
	}

	c := &character.Character{}

	// find character
	for _, data := range session.Characters {
		if data.Id == charId {
			c = data
			break
		}
	}

	// check if character exists
	if c.Id != charId {
		session.LogFatalf("Unable to find such character: %d", charId)
		return
	}

	// load additional character data
	req := character.DataReq{Id: c.Id}
	res := character.DataRes{}
	session.RPC.Call(rpc.LoadCharacterData, &req, &res)

	session.Character = c

	// serialize data
	eq, eqlen := c.Equipment.Serialize()
	inv, invlen := res.Inventory.Serialize()
	sk, sklen := res.Skills.Serialize()
	sl, sllen := res.Links.Serialize()

	pkt := network.NewWriter(CSCInitialize)
	pkt.WriteBytes(make([]byte, 57))
	pkt.WriteByte(0x00)
	pkt.WriteByte(0x14)
	pkt.WriteByte(session.ServerInstance.ChannelId)
	pkt.WriteBytes(make([]byte, 23))
	pkt.WriteByte(0xFF)
	pkt.WriteUint16(session.ServerConfig.MaxUsers)
	pkt.WriteUint32(0x8501A8C0)
	pkt.WriteUint16(0x985A)
	pkt.WriteInt32(0x01)
	pkt.WriteInt32(0x0100001F)

	pkt.WriteInt32(c.World)
	pkt.WriteInt32(0x00)
	pkt.WriteUint16(c.X)
	pkt.WriteUint16(c.Y)
	pkt.WriteUint64(c.Exp)
	pkt.WriteUint64(c.Alz)
	pkt.WriteUint64(c.WarExp)
	pkt.WriteUint32(c.Level)
	pkt.WriteInt32(0x00)

	pkt.WriteUint32(c.STR)
	pkt.WriteUint32(c.DEX)
	pkt.WriteUint32(c.INT)
	pkt.WriteUint32(c.PNT)
	pkt.WriteByte(c.MagicRank)
	pkt.WriteByte(c.SwordRank)
	pkt.WriteUint16(0x00) // padding for skillrank
	pkt.WriteUint32(0x00)
	pkt.WriteUint16(c.MaxHP)
	pkt.WriteUint16(c.CurrentHP)
	pkt.WriteUint16(c.MaxMP)
	pkt.WriteUint16(c.CurrentMP)
	pkt.WriteUint16(c.MaxSP)
	pkt.WriteUint16(c.CurrentSP)
	pkt.WriteUint16(0x00) //stats.DungeonPoints)
	pkt.WriteUint16(0x00)
	pkt.WriteInt32(0x2A30)
	pkt.WriteInt32(0x01)
	pkt.WriteUint16(0x00) //stats.SwordExp)
	pkt.WriteUint16(0x00) //stats.SwordPoint)
	pkt.WriteUint16(0x00) //stats.MagicExp)
	pkt.WriteUint16(0x00) //stats.MagicPoint)
	pkt.WriteUint16(0x00) //stats.SwordExpPoint)
	pkt.WriteUint16(0x00) //stats.MagicExpPoint)
	pkt.WriteInt32(0x00)
	pkt.WriteInt32(0x00)
	pkt.WriteInt32(0x00)  // honour pnt
	pkt.WriteUint64(0x00) // death penalty exp
	pkt.WriteUint64(0x00) // death hp
	pkt.WriteUint64(0x00) // death mp
	pkt.WriteUint16(0x00) // pk penalty // pk pna

	pkt.WriteUint32(0x8501A8C0) // chat ip
	pkt.WriteUint16(0x9858)     // chat port

	pkt.WriteUint32(0x8501A8C0) // ah ip
	pkt.WriteUint16(0x9859)     // ah port

	pkt.WriteByte(c.Nation)
	pkt.WriteInt32(0x00)
	pkt.WriteInt32(0x07) // warp code
	pkt.WriteInt32(0x07) // map code
	pkt.WriteUint32(c.Style.Get())
	pkt.WriteBytes(make([]byte, 39))

	pkt.WriteUint16(eqlen)
	pkt.WriteUint16(invlen)
	pkt.WriteUint16(sklen)
	pkt.WriteUint16(sllen)

	pkt.WriteBytes(make([]byte, 6))
	pkt.WriteUint16(0x00) // ap
	pkt.WriteUint32(0x00) // ap exp
	pkt.WriteInt16(0x00)
	pkt.WriteByte(0x00)   // blessing bead count
	pkt.WriteByte(0x00)   // active quest count
	pkt.WriteUint16(0x00) // period item count
	pkt.WriteBytes(make([]byte, 1023))

	pkt.WriteBytes(make([]byte, 128)) // quest dungeon flags
	pkt.WriteBytes(make([]byte, 128)) // mission dungeon flags

	pkt.WriteByte(0x00)              // Craft Lv 0
	pkt.WriteByte(0x00)              // Craft Lv 1
	pkt.WriteByte(0x00)              // Craft Lv 2
	pkt.WriteByte(0x00)              // Craft Lv 3
	pkt.WriteByte(0x00)              // Craft Lv 4
	pkt.WriteUint16(0x00)            // Craft Exp 0
	pkt.WriteUint16(0x00)            // Craft Exp 1
	pkt.WriteUint16(0x00)            // Craft Exp 2
	pkt.WriteUint16(0x00)            // Craft Exp 3
	pkt.WriteUint16(0x00)            // Craft Exp 4
	pkt.WriteBytes(make([]byte, 16)) // Craft Flags
	pkt.WriteUint32(0x00)            // Craft Type

	pkt.WriteInt32(0x10) // Help Window Index
	pkt.WriteBytes(make([]byte, 163))

	pkt.WriteUint32(0x00) // TotalPoints
	pkt.WriteUint32(0x00) // GeneralPoints
	pkt.WriteUint32(0x00) // QuestPoints
	pkt.WriteUint32(0x00) // DungeonPoints
	pkt.WriteUint32(0x00) // ItemPoints
	pkt.WriteUint32(0x00) // PVPPoints
	pkt.WriteUint32(0x00) // MissionWarPoints
	pkt.WriteUint32(0x00) // HuntingPoints
	pkt.WriteUint32(0x00) // CraftingPoints
	pkt.WriteUint32(0x00) // CommunityPoints
	pkt.WriteUint32(0x00) // SharedAchievments
	pkt.WriteUint32(0x00) // SpecialPoints

	pkt.WriteUint32(0x00)
	pkt.WriteUint32(0x00) // QuestsCount
	pkt.WriteUint32(0x00) // QuestFlagsCount
	pkt.WriteUint32(0x00)

	pkt.WriteByte(len(c.Name) + 1)
	pkt.WriteString(c.Name)

	pkt.WriteBytes(eq)
	pkt.WriteBytes(inv)
	pkt.WriteBytes(sk)
	pkt.WriteBytes(sl)

	// player is not moving anywhere, initialize begin/end movement variables
	c.BeginX = int16(c.X)
	c.BeginY = int16(c.Y)
	c.EndX = int16(c.X)
	c.EndY = int16(c.Y)

	// set-up inventory and links
	c.Inventory = &res.Inventory
	c.Links = &res.Links

	// set-up RPC and data inside inventory to sync with the database
	c.Inventory.Setup(session.RPC, c.Id)

	// set-up RPC and data inside equipment to sync with the database
	c.Equipment.Setup(session.RPC, c.Id)

	// set-up RPC and data inside links to sync with the database
	c.Links.Setup(session.RPC, c.Id)

	// prepare to enter the world
	session.World = session.WorldManager.FindWorld(c.World)
	if session.World == nil {
		session.LogFatalf("Unable to find world: %d for character %d",
			c.World, charId)
		return
	}

	session.SetState(StateInGame)

	session.Send(pkt)

	session.World.EnterWorld(session.SessionHandler)
	event.Trigger(event.PlayerJoin, session)
}

// UnInitialize Packet
func UnInitialize(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	_ = reader.ReadUint16() // index
	_ = reader.ReadByte()   // map id
	_ = reader.ReadByte()   // log out

	session.SetState(StateVerified)

	pkt := network.NewWriter(CSCUnInitialize)
	pkt.WriteByte(0) // result

	// complete - 0x00
	// fail - 0x01
	// ignored - 0x02
	// busy - 0x03
	// anti online game - 0x30

	session.Send(pkt)

	session.World.ExitWorld(session.SessionHandler, server.DelUserLogout)
}

// MessageEvent Packet
func MessageEvent(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	unk1 := reader.ReadInt16()
	msglen := reader.ReadInt16()
	_ = reader.ReadInt16()
	mtype := reader.ReadByte() // 0xA0 = normal; 0xA1 = trade; 0xA4 = roll dice
	msg := reader.ReadString(int(msglen) - 3)

	if strings.HasPrefix(msg, "#") {
		parts := strings.Split(msg, " ")
		command := parts[0][1:] // remove the '#'
		args := parts[1:]

		if err := script.ExecCommand(command, args, session); err != nil {
			log.Error("Failed to execute Lua command:", err)
		}

		return
	}

	pkt := network.NewWriter(NFYMessageEvent)
	pkt.WriteUint32(session.Character.Id)
	pkt.WriteByte(0) // 0x03 = [GM] prefix
	pkt.WriteByte(unk1)
	pkt.WriteByte(0)
	pkt.WriteByte(len(msg) + 3)
	pkt.WriteByte(0)
	pkt.WriteByte(254)
	pkt.WriteByte(254)
	pkt.WriteByte(mtype) // 0xA0 = normal; trade = 0xA1
	pkt.WriteString(msg)
	pkt.WriteByte(0)
	pkt.WriteByte(0)
	pkt.WriteByte(0)

	session.World.BroadcastSessionPacket(session.SessionHandler, pkt)
	time.Sleep(time.Second * 10)
}

// WarpCommand packet
func WarpCommand(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	warpId := reader.ReadByte()

	warp := session.World.FindWarp(warpId)
	if warp == nil {
		session.LogErrorf("Unable to find warp id: %d for character: %d",
			warpId, session.Character.Id)
		return
	}

	newWorld := session.WorldManager.FindWorld(warp.World)
	if newWorld == nil {
		session.LogErrorf("Unable to find new world by "+
			"warp id: %d for character: %d",
			warpId, session.Character.Id)
		return
	}

	pkt := network.NewWriter(CSCWarpCommand)
	pkt.WriteInt16(warp.Location[0].X) // pos x
	pkt.WriteInt16(warp.Location[0].Y) // pos y
	pkt.WriteInt32(0)                  // exp
	pkt.WriteInt32(0)                  // axp
	pkt.WriteInt32(0)                  // alz
	pkt.WriteInt32(0)                  // unk
	pkt.WriteInt16(session.GetUserIdx())
	pkt.WriteInt16(0x0100)
	pkt.WriteInt32(0x08)
	pkt.WriteByte(0)
	pkt.WriteInt32(warp.World)
	pkt.WriteInt32(0)
	pkt.WriteInt32(0)

	session.World.ExitWorld(session, server.DelUserWarp)

	x, y := warp.Location[0].X, warp.Location[0].Y

	session.Character.SetWorld(warp.World)
	session.Character.SetPosition(x, y)
	session.Character.SetMovement(x, y, x, y)

	session.Send(pkt)

	newWorld.EnterWorld(session)
}

func fillPlayerInfo(pkt *network.Writer, session *Session) {
	c := session.Character

	sx, sy, dx, dy := c.GetMovement()
	style, liveStyle := c.GetStyle()

	pkt.WriteUint32(c.Id)
	pkt.WriteUint32(session.GetUserIdx())
	pkt.WriteUint32(c.GetLevel())
	pkt.WriteInt32(0x01C2) // might be dwMoveBgnTime
	pkt.WriteUint16(sx)    // start
	pkt.WriteUint16(sy)
	pkt.WriteUint16(dx) // end
	pkt.WriteUint16(dy)
	pkt.WriteByte(0)
	pkt.WriteInt32(0)
	pkt.WriteInt16(0)
	pkt.WriteInt32(style.Get())
	pkt.WriteByte(liveStyle)
	pkt.WriteInt16(0)

	eq, eqlen := c.Equipment.SerializeEx()
	pkt.WriteInt16(eqlen)
	pkt.WriteInt16(0x00)

	for i := 0; i < 21; i++ {
		pkt.WriteByte(0)
	}

	pkt.WriteByte(len(c.Name) + 1)
	pkt.WriteString(c.Name)
	pkt.WriteByte(0) // guild name len
	// pkt.WriteString("guild name")

	pkt.WriteBytes(eq)
}

func NewUserSingle(session network.SessionHandler, reason server.NewUserType) *network.Writer {
	pkt := network.NewWriter(NFYNewUserList)
	pkt.WriteByte(1) // player num
	pkt.WriteByte(byte(reason))

	ses, ok := session.Retrieve().(*Session)
	if !ok {
		log.Error("Unable to parse client session!")
		return nil
	}

	fillPlayerInfo(pkt, ses)

	return pkt
}

func NewUserList(players map[uint16]network.SessionHandler, reason server.NewUserType) *network.Writer {
	online := len(players)

	pkt := network.NewWriter(NFYNewUserList)
	pkt.WriteByte(online)
	pkt.WriteByte(byte(reason))

	for _, v := range players {
		ses, ok := v.Retrieve().(*Session)
		if !ok {
			log.Error("Unable to parse client session!")
			return nil
		}

		fillPlayerInfo(pkt, ses)
	}

	return pkt
}

// DelUserList to all already connected players
func DelUserList(session network.SessionHandler, reason server.DelUserType) *network.Writer {
	ses, ok := session.Retrieve().(*Session)
	if !ok {
		log.Error("Unable to parse client session!")
		return nil
	}

	pkt := network.NewWriter(NFYDelUserList)
	pkt.WriteUint32(ses.Character.Id)
	pkt.WriteByte(byte(reason)) // type

	return pkt
}

// fixme
func SendMessage(session *Session, msg string) *network.Writer {
	pkt := network.NewWriter(NFYMessageEvent)
	pkt.WriteInt32(session.Character.Id)
	pkt.WriteByte(0) // 0x03 = [GM] prefix
	pkt.WriteByte(0x3F)
	pkt.WriteByte(0)
	pkt.WriteByte(len(msg) + 3)
	pkt.WriteByte(0)
	pkt.WriteByte(254)
	pkt.WriteByte(254)

	// normal = 0xA0;
	// trade  = 0xA1;
	// sys msg(right side) = 0xA3;
	// roll dice = 0xA4
	pkt.WriteByte(0xA4)
	pkt.WriteString(msg)
	pkt.WriteByte(0)
	pkt.WriteByte(0)
	pkt.WriteByte(0)

	return pkt
}
