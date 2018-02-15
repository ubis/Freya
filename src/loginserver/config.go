package main

import (
	"share/conf"
	"share/directory"
	"share/log"
)

// Config structure
type Config struct {
	// network
	Port int
	// client
	Version       int
	MagicKey      int
	CashWeb       string
	CashWebOdc    string
	CashWebCharge string
	GuildWeb      string
	Sns           string
	// master
	MasterIP   string
	MasterPort int
}

// Read attempts to open and parse server's configuration file
func (c *Config) Read() {
	log.Info("Reading configuration...")

	var location = directory.Root() + "/cfg/loginserver.ini"

	// parse configuration file...
	if err := conf.Open(location); err != nil {
		log.Fatal(err.Error())
		return
	}

	// read values from configuration...
	// network
	c.Port = conf.GetInt("network", "port", 38101)

	// client
	c.Version = conf.GetInt("client", "client_version", 0)
	c.MagicKey = conf.GetInt("client", "magic_key", 0)
	c.CashWeb = conf.GetString("client", "cashweb_url", "")
	c.CashWebOdc = conf.GetString("client", "cashweb_odc_url", "")
	c.CashWebCharge = conf.GetString("client", "cashweb_charge_url", "")
	c.GuildWeb = conf.GetString("client", "guildweb_url", "")
	c.Sns = conf.GetString("client", "sns_url", "")

	// master
	c.MasterIP = conf.GetString("master", "ip", "127.0.0.1")
	c.MasterPort = conf.GetInt("master", "port", 9001)
}
