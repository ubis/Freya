package def

import (
    "share/directory"
    "share/conf"
)

// Default values
const (
    C_Port     = 38101
    C_MaxUsers = 100
    C_Version  = 374
)

// Configuration struct
type Config struct {
    Port        int
    MaxUsers    int
    Version     int
}

// Attempts to read server configuration file
func (c *Config) Read() {
    log.Info("Reading configuration...")

    var location = directory.Root() + "/cfg/masterserver.ini"

    // parse configuration file...
    if err := conf.Open(location); err != nil {
        log.Fatalf("Couldn't read configuration file %s. %s", location, err.Error())
        return
    }

    // read values from configuration...
    c.Port     = conf.GetInt("network", "port", C_Port)
    c.MaxUsers = conf.GetInt("server", "max_users", C_MaxUsers)
    c.Version  = conf.GetInt("client", "client_version", C_Version)
}
