package server

import (
	"github.com/ubis/Freya/share/conf"
	"github.com/ubis/Freya/share/directory"
	"github.com/ubis/Freya/share/log"
)

type Config struct {
	Port int

	Version            int
	MagicKey           int
	IgnoreVersionCheck bool
	CashWeb_URL        string
	CashWeb_Odc_URL    string
	CashWeb_Charge_URL string
	GuildWeb_URL       string
	Sns_URL            string

	MasterIp   string
	MasterPort int

	ScriptDirectory string
}

// Attempts to read server configuration file
func (c *Config) Read() {
	log.Info("Reading configuration...")

	var location = directory.Root() + "/cfg/loginserver.ini"

	// parse configuration file...
	if err := conf.Open(location); err != nil {
		log.Fatal(err.Error())
	}

	// read values from configuration...
	c.Port = conf.GetInt("network", "port", 38101)

	c.Version = conf.GetInt("client", "client_version", 0)
	c.MagicKey = conf.GetInt("client", "magic_key", 0)
	c.IgnoreVersionCheck = conf.GetBool("client", "ignore_client_version", false)
	c.CashWeb_URL = conf.GetString("client", "cashweb_url", "")
	c.CashWeb_Odc_URL = conf.GetString("client", "cashweb_odc_url", "")
	c.CashWeb_Charge_URL = conf.GetString("client", "cashweb_charge_url", "")
	c.GuildWeb_URL = conf.GetString("client", "guildweb_url", "")
	c.Sns_URL = conf.GetString("client", "sns_url", "")

	c.MasterIp = conf.GetString("master", "ip", "1270.0.1")
	c.MasterPort = conf.GetInt("master", "port", 9001)

	c.ScriptDirectory = conf.GetString("script", "directory", "")
}
