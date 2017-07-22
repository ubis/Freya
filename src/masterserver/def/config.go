package def

import (
    "strconv"
    "fmt"
    "share/directory"
    "share/conf"
    "masterserver/database"
)

type Config struct {
    Port      int

    LoginIp   string
    LoginPort int
    LoginName string
    LoginUser string
    LoginPass string

    GameDB    map[int]*database.Database
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

    // load all game databases
    c.LoadGameDB()
}

// Returns LoginDatabase configuration string
func (c *Config) LoginDB() string {
    var str = c.LoginUser + ":" + c.LoginPass
    str += "@tcp(" + c.LoginIp + ":" + strconv.Itoa(c.LoginPort) + ")"
    str += "/" + c.LoginName + "?parseTime=true"
    return str
}

// Attemps to read all [1..255] GameDatabase configurations
func (c *Config) LoadGameDB() {
    c.GameDB = make(map[int]*database.Database)

    for i := 1; i < 256; i ++ {
        var section = fmt.Sprintf("game_%d", i)
        if conf.SectionExist(section) {
            c.GameDB[i] = &database.Database{
                conf.GetString(section, "ip", "127.0.0.1"),
                conf.GetInt(section, "port", 3306),
                conf.GetString(section, "database", "database"),
                conf.GetString(section, "username", "root"),
                conf.GetString(section, "password", ""),
                i,
                nil,
                "",
            }

            c.GameDB[i].Config = c.GetGameDB(c.GameDB[i])
        }
    }
}

// Returns GameDatabase configuration string
func (c *Config) GetGameDB(db *database.Database) string {
    var str = db.User + ":" + db.Pass
    str += "@tcp(" + db.Ip + ":" + strconv.Itoa(db.Port) + ")"
    str += "/" + db.Name + "?parseTime=true"
    return str
}