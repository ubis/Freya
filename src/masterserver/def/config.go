package def

import (
    "strconv"
    "fmt"
    "share/directory"
    "share/conf"
    "masterserver/database"
)

// Default values
const (
    C_Port    = 9001
    C_DB_IP   = "127.0.0.1"
    C_DB_PORT = 3306
    C_DB_NAME = "database"
    C_DB_USER = "root"
    C_DB_PASS = ""
)

// Configuration struct
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
        log.Fatalf("Couldn't read configuration file %s. %s", location, err.Error())
        return
    }

    // read values from configuration...
    c.Port = conf.GetInt("network", "port", C_Port)

    // login db
    c.LoginIp   = conf.GetString("login", "ip", C_DB_IP)
    c.LoginPort = conf.GetInt("login", "port", C_DB_PORT)
    c.LoginName = conf.GetString("login", "database", C_DB_NAME)
    c.LoginUser = conf.GetString("login", "username", C_DB_USER)
    c.LoginPass = conf.GetString("login", "password", C_DB_PASS)

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
                conf.GetString(section, "ip", C_DB_IP),
                conf.GetInt(section, "port", C_DB_PORT),
                conf.GetString(section, "database", C_DB_NAME),
                conf.GetString(section, "username", C_DB_USER),
                conf.GetString(section, "password", C_DB_PASS),
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