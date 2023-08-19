package script

import (
	"runtime"
	"runtime/debug"

	"github.com/ubis/Freya/share/event"
	"github.com/ubis/Freya/share/log"
	lua "github.com/yuin/gopher-lua"
)

// eventToLuaValue converts values from event.Event object to an equivalent
// slice of Lua values.
func eventToLuaValue(evt *event.Event, L *lua.LState) []lua.LValue {
	var luaValues []lua.LValue

	for {
		rawValue, ok := evt.Get()
		if !ok {
			break
		}

		switch v := rawValue.(type) {
		case string:
			luaValues = append(luaValues, lua.LString(v))
		case int:
			luaValues = append(luaValues, lua.LNumber(v))
		case float64:
			luaValues = append(luaValues, lua.LNumber(v))
		case bool:
			luaValues = append(luaValues, lua.LBool(v))
		default:
			ud := L.NewUserData()
			ud.Value = v
			luaValues = append(luaValues, ud)
		}
	}

	return luaValues
}

// Logging functions exported to Lua.

func infoMessage(L *lua.LState) int {
	log.Info(L.CheckString(1))
	return 0
}

func errorMessage(L *lua.LState) int {
	log.Error(L.CheckString(1))
	return 0
}

func debugMessage(L *lua.LState) int {
	log.Debug(L.CheckString(1))
	return 0
}

// reloadScripts reloads the list of scripts in the current Lua state.
func reloadScripts(L *lua.LState) int {
	log.Info("Reloading scripts...")

	for _, f := range scripts {
		log.Debugf("Reloading %s script...", f.file)

		for eventName, handler := range f.registeredEvents {
			event.Unregister(eventName, handler)
		}

		for key := range f.commandHandlers {
			delete(f.commandHandlers, key)
		}

		if err := L.DoFile(f.file); err != nil {
			log.Errorf("Error loading script %s: %v", f.file, err)
		}
	}

	return 0
}

// getGoVersion returns the Go runtime version to the Lua state.
func getGoVersion(L *lua.LState) int {
	L.Push(lua.LString(runtime.Version()))
	return 1
}

// getBuildInfo returns the build commit information for the application.
func getBuildInfo(L *lua.LState) int {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		L.Push(lua.LString("Unknown"))
		return 1
	}

	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			commit := setting.Value
			if len(commit) >= 6 {
				L.Push(lua.LString(commit[:6]))
				return 1
			}
			break
		}
	}

	L.Push(lua.LString("Unknown"))
	return 1
}

// addEventHandler registers a Lua function to handle a specific event.
func addEventHandler(L *lua.LState) int {
	eventName := L.CheckString(1)
	handler := L.CheckFunction(2)

	eventHandler := event.Register(eventName, func(e *event.Event) {
		lock.Lock()
		defer lock.Unlock()

		args := eventToLuaValue(e, L)

		if err := L.CallByParam(lua.P{
			Fn:      handler,
			NRet:    0,
			Protect: true,
		}, args...); err != nil {
			log.Error("Error calling Lua function: ", err)
		}
	})

	// store the registered handler
	for _, script := range scripts {
		if script.state == L {
			script.registeredEvents[eventName] = eventHandler
		}
	}

	return 0
}

// addCommandHandler associates a Lua function to handle a specific command.
func addCommandHandler(L *lua.LState) int {
	command := L.CheckString(1)
	handler := L.CheckFunction(2)

	for _, script := range scripts {
		if script.state == L {
			script.commandHandlers[command] = handler
			return 0
		}
	}

	return 0
}

// registerFunctions binds Go functions as global functions in the Lua state.
// This allows Lua scripts to invoke these Go functions directly.
func registerFunctions(L *lua.LState) {
	L.SetGlobal("infoMessage", L.NewFunction(infoMessage))
	L.SetGlobal("errorMessage", L.NewFunction(errorMessage))
	L.SetGlobal("debugMessage", L.NewFunction(debugMessage))

	L.SetGlobal("reloadScripts", L.NewFunction(reloadScripts))

	L.SetGlobal("getGoVersion", L.NewFunction(getGoVersion))
	L.SetGlobal("getBuildInfo", L.NewFunction(getBuildInfo))

	L.SetGlobal("addEventHandler", L.NewFunction(addEventHandler))
	L.SetGlobal("addCommandHandler", L.NewFunction(addCommandHandler))
}
