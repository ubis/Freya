package game

import (
	"os"
	"sync"

	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/share/directory"
	"github.com/ubis/Freya/share/log"
	"gopkg.in/yaml.v2"
)

// WorldManager manages world maps and their data
type WorldManager struct {
	mutex sync.RWMutex

	Worlds []*World
	Warps  []struct {
		World byte
		Warps []context.Warp
	}
}

// load reads data from a YAML file and deserializes it into the provided data structure.
func load(filename string, data any) error {
	s, err := os.ReadFile(directory.Root() + "data/" + filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(s, data)
}

// Initialize loads worlds, their data and initializes worlds.
func (wm *WorldManager) Initialize() {
	log.Info("Initializing World Manager...")

	// load world data
	if err := load("world.yml", &wm.Worlds); err != nil {
		log.Error("Failed to load world data:", err.Error())
		return
	}

	// load warps
	if err := load("warp.yml", &wm.Warps); err != nil {
		log.Error("Failed to load world warp data:", err.Error())
		return
	}

	log.Infof("Loaded %d world maps\n", len(wm.Worlds))
	log.Infof("Loaded %d warps\n", len(wm.Warps))

	// initialize each world
	for _, v := range wm.Worlds {
		v.Initialize(wm)
	}
}

// FindWorld finds a world by its ID.
func (wm *WorldManager) FindWorld(id byte) context.WorldHandler {
	wm.mutex.RLock()
	defer wm.mutex.RUnlock()

	for _, v := range wm.Worlds {
		if v.Id == id {
			return v
		}
	}

	return nil
}

// GetWarps returns a slice of warps for a specific world.
func (wm *WorldManager) GetWarps(world byte) []context.Warp {
	for _, v := range wm.Warps {
		if v.World != world {
			continue
		}

		return v.Warps
	}

	return nil
}
