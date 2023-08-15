package game

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/ubis/Freya/share/directory"
	"github.com/ubis/Freya/share/log"
	"gopkg.in/yaml.v2"
)

// dataDir constructs and returns the directory path for game data.
func dataDir() string {
	return directory.Root() + "data"
}

// load reads data from a YAML file and deserializes it into the provided data
// structure.
func load(filename string, data any) error {
	s, err := os.ReadFile(dataDir() + "/" + filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(s, data)
}

// loadMobs reads and returns a slice of Mobs associated with a specific World.
func (w *World) loadMobs() []*Mob {
	var mobs []*Mob

	err := load(fmt.Sprintf("world%d/mobs.yml", w.Id), &mobs)
	if err != nil {
		log.Error(err.Error())
	}

	return mobs
}

// loadThreadMap reads the binary file containing thread map information for a
// specific World.
func (w *World) loadThreadMap() {
	file := fmt.Sprintf("%s/world%d/tmap.bin", dataDir(), w.Id)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Error(err)
		return
	}

	f, err := os.Open(file)
	if err != nil {
		log.Error(err)
		return
	}

	defer f.Close()

	f.Seek(4, io.SeekStart) // iProcessNum

	var count int32
	binary.Read(f, binary.LittleEndian, &count)

	for i := int32(0); i < count; i++ {
		var column, row int32

		f.Seek(4, io.SeekCurrent) // iProcessIdx
		f.Seek(4, io.SeekCurrent) // iTileIdx
		binary.Read(f, binary.LittleEndian, &column)
		binary.Read(f, binary.LittleEndian, &row)
		f.Seek(4, io.SeekCurrent) // bIsEdge
		f.Seek(4, io.SeekCurrent) // pTileAttr

		// length: 5832
		// for j := 0; j < 9; j++ {
		// 	f.Seek(4*16, io.SeekCurrent)   // iTileIdxCurProcessLayer
		// 	f.Seek(4*9*16, io.SeekCurrent) // iTileIdxOthBoundary
		// 	f.Seek(4, io.SeekCurrent)      // iTileIdxOthProcessLayerNum
		// 	f.Seek(4, io.SeekCurrent)      // iTileIdxNum
		// }

		f.Seek(5832, io.SeekCurrent) // TileListEx

		cell := w.Grid[column][row]
		binary.Read(f, binary.LittleEndian, &cell.attribute) // TileAttr
	}
}
