package server

import (
	"github.com/ubis/Freya/share/conf"
	"github.com/ubis/Freya/share/directory"
	"github.com/ubis/Freya/share/log"
)

type Config struct {
	PublicIp   string
	Port       int
	MaxUsers   int
	UseLocalIp bool

	ServerType        int
	IgnoreSubPassword bool

	MasterIp   string
	MasterPort int

	ScriptDirectory string
}

// Attempts to read server configuration file
func (c *Config) Read(name string) {
	log.Info("Reading configuration...")

	var location = directory.Root() + "/cfg/" + name + ".ini"

	// parse configuration file...
	if err := conf.Open(location); err != nil {
		log.Fatal(err.Error())
		return
	}

	// read values from configuration...
	c.PublicIp = conf.GetString("network", "ip", "127.0.0.1")
	c.Port = conf.GetInt("network", "port", 38111)
	c.MaxUsers = conf.GetInt("network", "max_users", 100)
	c.UseLocalIp = conf.GetBool("network", "use_local_ip", false)

	c.ServerType = conf.GetInt("server", "server_type", 0)
	c.IgnoreSubPassword = conf.GetBool("server", "ignore_sub_password", false)

	c.MasterIp = conf.GetString("master", "ip", "127.0.0.1")
	c.MasterPort = conf.GetInt("master", "port", 9001)

	c.ScriptDirectory = conf.GetString("script", "directory", "")
}
