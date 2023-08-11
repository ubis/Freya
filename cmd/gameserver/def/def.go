package def

import (
	"os"
	"strconv"

	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/rpc"
)

var ServerConfig = &Config{}
var ServerSettings = &Settings{}
var NetworkManager = &network.Network{}
var PacketHandler = &network.PacketHandler{}
var RPCHandler = &rpc.Client{}

// init function, which runs before main()
func init() {
	if len(os.Args) > 2 {
		if id, err := strconv.Atoi(os.Args[1]); err == nil {
			ServerSettings.ServerId = id
		} else {
			ServerSettings.ServerId = 1
		}

		if id, err := strconv.Atoi(os.Args[2]); err == nil {
			ServerSettings.ChannelId = id
		} else {
			ServerSettings.ChannelId = 1
		}
	} else {
		ServerSettings.ServerId = 1
		ServerSettings.ChannelId = 1
	}
}

// Returns GameServer name with id's
func GetName() string {
	var str = "GameServer_" + strconv.Itoa(ServerSettings.ServerId)
	str += "_" + strconv.Itoa(ServerSettings.ChannelId)

	return str
}
