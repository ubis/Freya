package def

import (
    "share/directory"
    "share/conf"
)

// Default values
const (
    C_Port          = 38101

    C_Version       = 0
    C_MagicKey      = 0

    C_MasterIp      = "127.0.0.1"
    C_MasterPort    = 9001
)

// Configuration struct
type Config struct {
    Port               int

    Version            int
    MagicKey           int
    CashWeb_URL        string
    CashWeb_Odc_URL    string
    CashWeb_Charge_URL string
    GuildWeb_URL       string
    Sns_URL            string

    MasterIp           string
    MasterPort         int
}

// Attempts to read server configuration file
func (c *Config) Read() {
    log.Info("Reading configuration...")

    var location = directory.Root() + "/cfg/loginserver.ini"

    // parse configuration file...
    if err := conf.Open(location); err != nil {
        log.Fatalf("Couldn't read configuration file %s. %s", location, err.Error())
        return
    }

    // read values from configuration...
    c.Port = conf.GetInt("network", "port", C_Port)

    c.Version            = conf.GetInt("client", "client_version", C_Version)
    c.MagicKey           = conf.GetInt("client", "magic_key", C_MagicKey)
    c.CashWeb_URL        = conf.GetString("client", "cashweb_url", "")
    c.CashWeb_Odc_URL    = conf.GetString("client", "cashweb_odc_url", "")
    c.CashWeb_Charge_URL = conf.GetString("client", "cashweb_charge_url", "")
    c.GuildWeb_URL       = conf.GetString("client", "guildweb_url", "")
    c.Sns_URL            = conf.GetString("client", "sns_url", "")

    c.MasterIp   = conf.GetString("master", "ip", C_MasterIp)
    c.MasterPort = conf.GetInt("master", "port", C_MasterPort)
}