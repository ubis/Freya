package packet

import (
	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/models/inventory"
	"github.com/ubis/Freya/share/network"
	lua "github.com/yuin/gopher-lua"
)

type sessionPacketFunc struct{}

type clientMessageFunc struct {
	Fn func(*network.Session, string) *network.Writer
}

type playerGetLevelFunc struct {
	Fn func(*network.Session) int
}

type playerSetLevelFunc struct {
	Fn func(*network.Session, int)
}

type playerPositionFunc struct{}
type playerDropItemFunc struct{}

func (cmf sessionPacketFunc) Call(L *lua.LState) []lua.LValue {
	ud := L.CheckUserData(1) // session
	num := L.CheckNumber(2)  // opcode
	tbl := L.CheckTable(3)   // byte array

	session, ok := ud.Value.(*network.Session)
	if !ok {
		return nil
	}

	data := make([]byte, tbl.Len())
	for i := 1; i <= tbl.Len(); i++ {
		val := tbl.RawGetInt(i)
		if byteVal, ok := val.(lua.LNumber); ok {
			data[i-1] = byte(byteVal)
		} else {
			L.ArgError(2, "Expected byte array")
			return nil
		}
	}

	pkt := network.NewWriter(uint16(num))
	pkt.WriteBytes(data)

	session.Send(pkt)

	return nil
}

func (cmf clientMessageFunc) Call(L *lua.LState) []lua.LValue {
	ud := L.CheckUserData(1)
	msg := L.CheckString(2)

	session, ok := ud.Value.(*network.Session)
	if !ok {
		return nil
	}

	pkt := cmf.Fn(session, msg)
	session.Send(pkt)

	return nil
}

func (cmf playerGetLevelFunc) Call(L *lua.LState) []lua.LValue {
	ud := L.CheckUserData(1)

	session, ok := ud.Value.(*network.Session)
	if !ok {
		return nil
	}

	level := cmf.Fn(session)

	return []lua.LValue{lua.LNumber(level)}
}

func (cmf playerSetLevelFunc) Call(L *lua.LState) []lua.LValue {
	ud := L.CheckUserData(1)
	level := L.CheckInt(2)

	session, ok := ud.Value.(*network.Session)
	if !ok {
		return nil
	}

	cmf.Fn(session, level)

	return nil
}

func (cmd playerPositionFunc) Call(L *lua.LState) []lua.LValue {
	ud := L.CheckUserData(1)
	session, ok := ud.Value.(*network.Session)
	if !ok {
		return nil
	}

	var x, y byte

	if err := context.GetCharPosition(session, &x, &y); err != nil {
		log.Error(err.Error())
		return nil
	}

	return []lua.LValue{lua.LNumber(x), lua.LNumber(y)}
}

func (cmd playerDropItemFunc) Call(L *lua.LState) []lua.LValue {
	ud := L.CheckUserData(1)
	kind := L.CheckNumber(2)
	option := L.CheckNumber(3)

	session, ok := ud.Value.(*network.Session)
	if !ok {
		return nil
	}

	ctx, err := context.Parse(session)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	ctx.Mutex.RLock()
	world := ctx.World
	id := ctx.Char.Id
	x := int(ctx.Char.X)
	y := int(ctx.Char.Y)
	ctx.Mutex.RUnlock()

	item := &inventory.Item{
		Kind:   uint32(kind),
		Option: int32(option),
	}

	if !world.DropItem(item, id, x, y) {
		log.Error("unable to drop new item", item)
	}

	return nil
}
