package packet

import (
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
