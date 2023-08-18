package packet

import (
	"github.com/ubis/Freya/share/network"
	lua "github.com/yuin/gopher-lua"
)

type clientMessageFunc struct {
	Fn func(string) *network.Writer
}

func (cmf clientMessageFunc) Call(L *lua.LState) {
	ud := L.CheckUserData(1)
	msg := L.CheckString(2)

	session, ok := ud.Value.(*network.Session)
	if !ok {
		return
	}

	pkt := cmf.Fn(msg)
	if pkt != nil {
		session.Send(pkt)
	}
}
