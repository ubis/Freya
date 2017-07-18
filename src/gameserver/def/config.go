package def

import (
    "share/directory"
    "share/conf"
)

// Default values
const (
    C_PublicIp   = "127.0.0.1"
    C_Port       = 38111
    C_MaxUsers   = 100

    C_ServerType = 0

    C_MasterIp   = "127.0.0.1"
    C_MasterPort = 9001
)

// Configuration struct
type Config struct {
    PublicIp   string
    Port       int
    MaxUsers   int

    ServerType int

    MasterIp   string
    MasterPort int
}

// Attempts to read server configuration file
func (c *Config) Read() {
    log.Info("Reading configuration...")

    var location = directory.Root() + "/cfg/" + GetName() + ".ini"

    // parse configuration file...
    if err := conf.Open(location); err != nil {
        log.Fatalf("Couldn't read configuration file %s. %s", location, err.Error())
        return
    }

    // read values from configuration...
    c.PublicIp = conf.GetString("network", "ip", C_PublicIp)
    c.Port     = conf.GetInt("network", "port", C_Port)
    c.MaxUsers = conf.GetInt("network", "max_users", C_MaxUsers)

    c.ServerType = conf.GetInt("server", "server_type", C_ServerType)

    c.MasterIp   = conf.GetString("master", "ip", C_MasterIp)
    c.MasterPort = conf.GetInt("master", "port", C_MasterPort)
}