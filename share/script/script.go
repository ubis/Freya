package script

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ubis/Freya/share/event"
	"github.com/ubis/Freya/share/log"
	lua "github.com/yuin/gopher-lua"
)

// LuaCallable represents an interface for objects that can be
// called within a Lua state.
type LuaCallable interface {
	Call(L *lua.LState) []lua.LValue
}

// script encapsulates a single loaded Lua script, including its state,
// file path, and any registered command handlers.
type script struct {
	state            *lua.LState
	file             string
	commandHandlers  map[string]*lua.LFunction
	registeredEvents map[string]event.HandlerId
}

// scripts is a global slice containing all loaded Lua scripts.
var scripts []script
var lock sync.Mutex

// loadScript loads a Lua script from the provided file path.
// It initializes a new Lua state for the script and registers
// available functions.
func loadScript(file string) {
	L := lua.NewState()

	scripts = append(scripts, script{
		state:            L,
		file:             file,
		commandHandlers:  make(map[string]*lua.LFunction),
		registeredEvents: make(map[string]event.HandlerId),
	})

	registerFunctions(L)

	log.Debugf("Loading %s script...", file)

	if err := L.DoFile(file); err != nil {
		log.Errorf("Error loading script %s: %v", file, err)
	}
}

// Initialize loads all Lua scripts from the provided directory.
func Initialize(directory string) {
	if len(directory) == 0 {
		log.Warning("Script directory path is empty!")
		return
	}

	log.Info("Initializing scripting engine...")
	log.Infof("Loading scripts from %s directory...", directory)

	files, err := os.ReadDir(directory)
	if err != nil {
		log.Error(err)
		return
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".lua") {
			loadScript(filepath.Join(directory, f.Name()))
		}
	}
}

// RegisterFunc registers a Go function as a global function accessible from
// all Lua scripts.
func RegisterFunc(funcName string, callable LuaCallable) {
	for _, script := range scripts {
		L1 := script.state

		L1.SetGlobal(funcName, L1.NewFunction(func(L *lua.LState) int {
			args := callable.Call(L)

			for _, val := range args {
				L.Push(val)
			}

			return len(args)
		}))
	}
}

// ExecCommand executes a Lua command using the provided arguments and session.
// It searches for the command handler across all loaded scripts and
// executes the first match found.
func ExecCommand(command string, args []string, session any) error {
	var fn *lua.LFunction
	var file *script

	lock.Lock()
	defer lock.Unlock()

	for _, script := range scripts {
		handler, exists := script.commandHandlers[command]
		if exists {
			fn = handler
			file = &script
			break
		}
	}

	if fn == nil || file == nil {
		return fmt.Errorf("command %s not found", command)
	}

	L := file.state

	ud := L.NewUserData()
	ud.Value = session
	L.SetMetatable(ud, L.GetTypeMetatable("session_ud"))

	luaArgs := make([]lua.LValue, len(args)+1)
	luaArgs[0] = ud
	for i, arg := range args {
		luaArgs[i+1] = lua.LString(arg)
	}

	err := L.CallByParam(lua.P{
		Fn:      fn,
		NRet:    0,
		Protect: true,
	}, luaArgs...)

	return err
}
