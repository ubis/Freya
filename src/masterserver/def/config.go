package def

import (
    "strconv"
    "share/directory"
    "share/conf"
)

// Default values
const (
    C_Port    = 9001
    C_DB_IP   = "127.0.0.1"
    C_DB_PORT = 3306
    C_DB_NAME = "database"
)

// Configuration struct
type Config struct {
    Port      int

    LoginIp   string
    LoginPort int
    LoginName string
    LoginUser string
    LoginPass string
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
    c.Port = conf.GetInt("network", "port", C_Port)

    // login db
    c.LoginIp   = conf.GetString("login", "ip", C_DB_IP)
    c.LoginPort = conf.GetInt("login", "port", C_DB_PORT)
    c.LoginName = conf.GetString("login", "database", C_DB_NAME)
    c.LoginUser = conf.GetString("login", "username", "root")
    c.LoginPass = conf.GetString("login", "password", "")
}

func (c *Config) LoginDB() string {
    str := c.LoginUser + ":" + c.LoginPass
    str += "@tcp(" + c.LoginIp + ":" + strconv.Itoa(c.LoginPort) + ")"
    str += "/" + c.LoginName + "?parseTime=true"
    return str
}