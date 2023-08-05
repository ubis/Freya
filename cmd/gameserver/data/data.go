package data

import (
	"os"

	"github.com/ubis/Freya/share/directory"
	"github.com/ubis/Freya/share/log"
	"gopkg.in/yaml.v2"
)

type Loader struct {
	*WarpData
}

// Initializes DataLoader
func (dl *Loader) Init() {
	log.Info("Loading data...")

	dl.WarpData = &WarpData{}
	dl.load("warp.yml", dl.WarpData)
}

// Deserializes data from file to specified struct
func (dl *Loader) load(filename string, data interface{}) {
	var s, err = os.ReadFile(directory.Root() + "/data/" + filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = yaml.Unmarshal(s, data)
	if err != nil {
		log.Fatal(err.Error())
	}
}
