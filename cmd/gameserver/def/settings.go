package def

import "github.com/ubis/Freya/share/models/server"

type Settings struct {
	server.Settings
	ServerId  int
	ChannelId int
}
