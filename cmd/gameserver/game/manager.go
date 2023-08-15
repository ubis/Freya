package game

import (
	"github.com/ubis/Freya/cmd/gameserver/context"
	"github.com/ubis/Freya/share/log"
)

// WorldManager manages world maps and their data
type WorldManager struct {
	Worlds []*World
	Warps  []struct {
		World byte
		Warps []context.Warp
	}
	Mobs []*Mob
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

	// load mobs
	if err := load("mobs.yml", &wm.Mobs); err != nil {
		log.Error("Failed to load world mob data:", err.Error())
		return
	}

	log.Infof("Loaded %d world maps\n", len(wm.Worlds))
	log.Infof("Loaded %d warps\n", len(wm.Warps))
	log.Infof("Loaded %d mobs\n", len(wm.Mobs))

	// initialize each world
	for _, v := range wm.Worlds {
		v.Initialize(wm)
	}
}

// FindWorld finds a world by its ID.
func (wm *WorldManager) FindWorld(id byte) context.WorldHandler {
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

// GetMob returns the mob that matches a given species ID.
func (wm *WorldManager) GetMob(species int) *Mob {
	for _, v := range wm.Mobs {
		if v.Species != species {
			continue
		}

		return v
	}

	return nil
}
