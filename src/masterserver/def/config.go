package def

import (
    "strconv"
    "share/directory"
    "share/conf"
)

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
        log.Fatal(err.Error())
        return
    }

    // read values from configuration...
    c.Port = conf.GetInt("network", "port", 9001)

    // login db
    c.LoginIp   = conf.GetString("login", "ip", "127.0.0.1")
    c.LoginPort = conf.GetInt("login", "port", 3306)
    c.LoginName = conf.GetString("login", "database", "database")
    c.LoginUser = conf.GetString("login", "username", "root")
    c.LoginPass = conf.GetString("login", "password", "")
}

// Returns LoginDatabase configuration string
func (c *Config) LoginDB() string {
    str := c.LoginUser + ":" + c.LoginPass
    str += "@tcp(" + c.LoginIp + ":" + strconv.Itoa(c.LoginPort) + ")"
    str += "/" + c.LoginName
    return str
}