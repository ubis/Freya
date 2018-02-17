package main

import (
	"gameserver/net"
	"share/conf"
	"share/directory"
	"share/log"
	"strconv"
)

// Config structure
type Config struct {
	// internal
	ServerID int
	GroupID  int

	// network
	PublicIp string
	Port     int
	MaxUsers int
	// server
	ServerType int
	// master
	MasterIP   string
	MasterPort int
}

// Read attempts to open and parse server's configuration file
func (c *Config) Read() {
	log.Info("Reading configuration...")

	location := directory.Root() + "/cfg/" + c.GetName() + ".ini"

	// parse configuration file...
	if err := conf.Open(location); err != nil {
		log.Fatal(err.Error())
		return
	}

	// read values from configuration...
	// network
	c.PublicIp = conf.GetString("network", "ip", "127.0.0.1")
	c.Port = conf.GetInt("network", "port", 38111)
	c.MaxUsers = conf.GetInt("network", "max_users", 100)

	// server
	c.ServerType = conf.GetInt("server", "server_type", 0)

	// master
	c.MasterIP = conf.GetString("master", "ip", "127.0.0.1")
	c.MasterPort = conf.GetInt("master", "port", 9001)
}

// Assign configuration for Packet structure
func (c *Config) Assign(p *net.Packet) {
	p.ServerID = c.ServerID
	p.GroupID = c.GroupID
	p.ServerType = c.ServerType
}

// GetName returns server's name with sid/gid
func (c *Config) GetName() string {
	str := "GameServer_" + strconv.Itoa(c.ServerID)
	str += "_" + strconv.Itoa(c.GroupID)

	return str
}
