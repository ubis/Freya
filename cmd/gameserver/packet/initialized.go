package packet

import (
	"github.com/ubis/Freya/share/models/character"
	"github.com/ubis/Freya/share/models/server"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"

	"github.com/ubis/Freya/share/log"
)

// Initialized Packet
func Initialized(session *network.Session, reader *network.Reader) {
	charId := reader.ReadInt32()

	if !session.Data.Verified || !session.Data.LoggedIn || session.DataEx == nil {
		log.Errorf("User is not verified (char: %d)", charId)
		return
	}

	ctx, ok := session.DataEx.(*context)
	if !ok {
		log.Error("Unable to retrieve user context (id: %d)",
			session.Data.AccountId)
		return
	}

	// verify char id
	if (charId >> 3) != session.Data.AccountId {
		log.Errorf("User is using invalid character id (id: %d, char: %d)",
			session.Data.AccountId, charId)
		return
	}

	c := character.Character{}

	if len(session.Data.CharacterList) == 0 {
		// fetch characters
		reqList := character.ListReq{
			Account: session.Data.AccountId,
			Server:  byte(g_ServerSettings.ServerId),
		}
		resList := character.ListRes{}
		g_RPCHandler.Call(rpc.LoadCharacters, &reqList, &resList)

		session.Data.CharacterList = resList.List
	}

	// fetch character
	for _, data := range session.Data.CharacterList {
		if data.Id == charId {
			c = data
			break
		}
	}

	// check if character exists
	if c.Id != charId {
		log.Errorf("User is using invalid character id (id: %d, char: %d)",
			session.Data.AccountId, charId)
		return
	}

	// load additional character data
	req := character.DataReq{
		Server: byte(g_ServerSettings.ServerId),
		Id:     c.Id,
	}
	res := character.DataRes{}
	g_RPCHandler.Call(rpc.LoadCharacterData, req, &res)

	// serialize data
	eq, eqlen := c.Equipment.Serialize()
	inv, invlen := res.Inventory.Serialize()
	sk, sklen := res.Skills.Serialize()
	sl, sllen := res.Links.Serialize()

	pkt := network.NewWriter(INITIALIZED)
	pkt.WriteBytes(make([]byte, 57))
	pkt.WriteByte(0x00)
	pkt.WriteByte(0x14)
	pkt.WriteByte(g_ServerSettings.ChannelId)
	pkt.WriteBytes(make([]byte, 23))
	pkt.WriteByte(0xFF)
	pkt.WriteUint16(g_ServerConfig.MaxUsers)
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
	pkt.WriteByte(c.SwordRank)
	pkt.WriteByte(c.MagicRank)
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

	pkt.WriteInt32(0x00) // Help Window Index
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

	session.Send(pkt)

	// player is not moving anywhere, initialize begin/end movement variables
	c.BeginX = int16(c.X)
	c.BeginY = int16(c.Y)
	c.EndX = int16(c.X)
	c.EndY = int16(c.Y)

	c.Inventory = res.Inventory

	ctx.mutex.Lock()
	ctx.char = &c
	ctx.mutex.Unlock()

	notifyNewPlayer(session)     // notify others about new player
	notifyOnlinePlayers(session) // notify all players to our new player
}

// Uninitialze Packet
func Uninitialze(session *network.Session, reader *network.Reader) {
	_ = reader.ReadUint16() // index
	_ = reader.ReadByte()   // map id
	_ = reader.ReadByte()   // log out

	pkt := network.NewWriter(UNINITIALZE)
	pkt.WriteByte(0) // result

	// complete - 0x00
	// fail - 0x01
	// ignored - 0x02
	// busy - 0x03
	// anti online game - 0x30

	session.Send(pkt)

	DelUserList(session, server.DelUserLogout) // notify that player is gone
}

// notifyNewPlayer to all already connected players
func notifyNewPlayer(session *network.Session) {
	pkt := network.NewWriter(NEWUSERLIST)
	pkt.WriteUint16(1) // player num

	g_NetworkManager.RLock()
	fillPlayerInfo(pkt, session)
	g_NetworkManager.RUnlock()

	g_NetworkManager.SendToAllExcept(pkt, session)
}

// notifyOnlinePlayers to newly connected player
func notifyOnlinePlayers(session *network.Session) {
	online := g_NetworkManager.GetOnlineUsers()
	if online == 1 {
		// no point, we are alone
		return
	}

	online--
	sessions := g_NetworkManager.GetSessions()

	pkt := network.NewWriter(NEWUSERLIST)
	pkt.WriteUint16(online)

	g_NetworkManager.RLock()
	for _, v := range sessions {
		if v == session || v.DataEx == nil {
			// we don't want to send player data to our newly connected session
			// or invalid clients - connected but not sent Connect2Svr packet
			continue
		}

		if _, ok := v.DataEx.(*context); !ok {
			// this is ok - we might have users in the lobby and not in-game
			// if user is not in-game, DataEx.context will be nil
			continue
		}

		fillPlayerInfo(pkt, v)

	}
	g_NetworkManager.RUnlock()

	session.Send(pkt)
}

// fillPlayerInfo with player character data
func fillPlayerInfo(pkt *network.Writer, session *network.Session) {
	ctx, ok := session.DataEx.(*context)
	if !ok {
		log.Error("Unable to retrieve user context (id: %d)",
			session.Data.AccountId)
		return
	}

	ctx.mutex.RLock()
	defer ctx.mutex.RUnlock()

	c := ctx.char

	if c == nil {
		// client is not ready
		return
	}

	pkt.WriteUint32(c.Id)
	pkt.WriteUint32(session.UserIdx)
	pkt.WriteUint32(c.Level)
	pkt.WriteInt32(0x01C2)    // might be dwMoveBgnTime
	pkt.WriteUint16(c.BeginX) // start
	pkt.WriteUint16(c.BeginY)
	pkt.WriteUint16(c.EndX) // end
	pkt.WriteUint16(c.EndY)
	pkt.WriteByte(0)
	pkt.WriteInt32(0)
	pkt.WriteInt16(0)
	pkt.WriteInt32(c.Style.Get())
	pkt.WriteByte(0) // animation id aka "live style"
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

// DelUserList to all already connected players
func DelUserList(session *network.Session, reason server.DelUserType) {
	if session.DataEx == nil {
		// we might have invalid clients, ignore
		return
	}

	ctx, ok := session.DataEx.(*context)
	if !ok {
		// we might have invalid clients, ignore
		return
	}

	ctx.mutex.RLock()
	if ctx.char == nil {
		// user might be (or was) in the lobby
		ctx.mutex.RUnlock()
		return
	}

	charId := ctx.char.Id

	if reason != server.DelUserWarp {
		ctx.char = nil // we are no longer in the world
	}

	ctx.mutex.RUnlock()

	pkt := network.NewWriter(DELUSERLIST)
	pkt.WriteUint32(charId)
	pkt.WriteByte(byte(reason)) // type

	/* types:
	 * dead = 0x10
	 * warp = 0x11
	 * logout = 0x12
	 * retn = 0x13
	 * dissapear = 0x14
	 * nfsdead = 0x15
	 */

	g_NetworkManager.SendToAllExcept(pkt, session)
}

// MessageEvnt Packet
func MessageEvnt(session *network.Session, reader *network.Reader) {
	unk1 := reader.ReadInt16()
	msglen := reader.ReadInt16()
	_ = reader.ReadInt16()
	mtype := reader.ReadByte() // 0xA0 = normal; 0xA1 = trade
	msg := reader.ReadString(int(msglen))

	if session.DataEx == nil {
		// we have invalid user, ignore
		return
	}

	ctx, ok := session.DataEx.(*context)
	if !ok {
		// we have invalid user, ignore
		return
	}

	ctx.mutex.RLock()
	if ctx.char == nil {
		// user is in the lobby, we cannot receive such messages
		ctx.mutex.RUnlock()
		return
	}

	charId := ctx.char.Id
	ctx.mutex.RUnlock()

	pkt := network.NewWriter(NFY_MESSAGEEVNT)
	pkt.WriteUint32(charId)
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

	g_NetworkManager.SendToAll(pkt)
}

// WarpCommand packet
func WarpCommand(session *network.Session, reader *network.Reader) {
	warpId := reader.ReadByte()

	ctx, err := parseSessionContext(session)
	if err != nil {
		log.Error(err.Error())
		return
	}

	ctx.mutex.RLock()
	world := ctx.char.World
	ctx.mutex.RUnlock()

	warp := g_DataLoader.WarpData.FindWarp(world, warpId)
	if warp == nil {
		log.Error("Unable to find warp:", warpId)
		return
	}

	pkt := network.NewWriter(WARPCOMMAND)
	pkt.WriteInt16(warp.Location[0].X) // pos x
	pkt.WriteInt16(warp.Location[0].Y) // pos y
	pkt.WriteInt32(0)                  // exp
	pkt.WriteInt32(0)                  // axp
	pkt.WriteInt32(0)                  // alz
	pkt.WriteInt32(0)                  // unk
	pkt.WriteInt16(session.UserIdx)
	pkt.WriteInt16(0x0100)
	pkt.WriteInt32(0x08)
	pkt.WriteByte(0)
	pkt.WriteInt32(warp.World)
	pkt.WriteInt32(0)
	pkt.WriteInt32(0)

	DelUserList(session, server.DelUserWarp)

	session.Send(pkt)

	ctx.mutex.Lock()
	ctx.char.World = byte(warp.World)
	ctx.char.X = byte(warp.Location[0].X)
	ctx.char.Y = byte(warp.Location[0].Y)
	ctx.char.BeginX = int16(warp.Location[0].X)
	ctx.char.BeginY = int16(warp.Location[0].Y)
	ctx.char.EndX = int16(warp.Location[0].X)
	ctx.char.EndY = int16(warp.Location[0].Y)
	ctx.mutex.Unlock()

	// this will produce little glitches - player re-appears, though it is in
	// the different world
	// this is because we don't have proper map management (yet)
	notifyNewPlayer(session)     // notify others about new player
	notifyOnlinePlayers(session) // notify all players to our new player
}
