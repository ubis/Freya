package packet

import (
	"github.com/ubis/Freya/share/network"
	lua "github.com/yuin/gopher-lua"
)

type clientMessageFunc struct {
	Fn func(*network.Session, string) *network.Writer
}

type playerGetLevelFunc struct {
	Fn func(*network.Session) int
}

type playerSetLevelFunc struct {
	Fn func(*network.Session, int)
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
