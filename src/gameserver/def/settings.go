package def

import "share/models/server"

type Settings struct {
	server.Settings
	ServerId  int
	ChannelId int
}
