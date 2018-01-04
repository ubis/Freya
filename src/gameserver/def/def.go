package def

import (
	"github.com/op/go-logging"
	"os"
	"share/logger"
	"share/network"
	"share/rpc"
	"strconv"
)

var log *logging.Logger

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

	log = logger.Init(GetName())
}

// Returns GameServer name with id's
func GetName() string {
	var str = "gameserver_" + strconv.Itoa(ServerSettings.ServerId)
	str += "_" + strconv.Itoa(ServerSettings.ChannelId)

	return str
}
