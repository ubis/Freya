package main

import (
	"strconv"
	"share/lib/goini"
	"share/directory"
)

// Default values
const (
	C_Port 	 	= 38101
	C_MaxUsers 	= 100
)

// Configuration struct
type Config struct {
	Port		int
	MaxUsers	int
}

// Attempts to read server configuration file
func (c *Config) Read() {
	log.Info("Reading configuration...")

	var location = directory.Root() + "/cfg/loginserver.ini"
	var ini 	 = goini.New()

	// parse configuration file...
	if err := ini.ParseFile(location); err != nil {
		log.Fatalf("Couldn't read configuration file %s. %s", location, err.Error())
		return
	}

	// read values from configuration...
	if port, err := ini.SectionGet("network", "port"); err != true {
		c.Port = C_Port
	} else {
		c.Port, _ = strconv.Atoi(port)
	}

	if maxUsers, err := ini.SectionGet("server", "max_users"); err != true {
		c.MaxUsers = C_MaxUsers
	} else {
		c.MaxUsers, _ = strconv.Atoi(maxUsers)
	}
}