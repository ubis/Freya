package conf

import (
    "io/ioutil"
    "strings"
    "errors"
    "strconv"
    "share/logger"
)

type keyValue map[string]string
type section  map[string]keyValue

var log = logger.Instance()

var sections = make(section)

// Attempts to open and read configuration file. Error is returned on fail
func Open(path string) error {
    var f, err = ioutil.ReadFile(path)
    if err != nil {
        return err
    }

    var lines   = strings.Split(string(f), "\n")
    var section = ""

    for i := 0; i < len(lines); i++ {
        var line = strings.TrimSpace(lines[i])
        line     = strings.Replace(line, "\r", "", -1)
        var size = len(line)

        // empty line
        if size == 0 {
            continue
        }

        // comment
        if line[0] == ';' || line[0] == '#' {
            continue
        }

        // section
        if line[0] == '[' {
            section = strings.Split(line, "]")[0]
            section = strings.Replace(section,"[","", 1)
            section = strings.TrimSpace(section)

            if len(section) == 0 {
                return errors.New("Error parsing configuration file!")
            }
        } else {
            // value and key
            var data  = strings.Split(line, ";")[0]
            var split = strings.Index(line, "=")
            var key   = strings.TrimSpace(data[:split])
            var value = strings.TrimSpace(data[split + 1:])

            // create section if it wasn't
            if sections[section] == nil {
                sections[section] = make(keyValue)
            }

            var tmp = sections[section]
            tmp[key] = value

            // don't show password in plain-text
            if strings.ToLower(key) == "password" {
                value = "******"
            }
        }
    }

    return nil
}

// Returns string value from configuration file. If section or key isn't found,
// default value will be returned
func GetString(section string, key string, def string) string {
    if value, err := get(section, key); err == nil {
        return value
    }

    return def
}

// Returns int value from configuration file. If section or key isn't found,
// default value will be returned
func GetInt(section string, key string, def int) int {
    if value, err := get(section, key); err == nil {
        var tmp, _ = strconv.Atoi(value)
        return tmp
    }

    return def
}

// Returns bool value from configuration file. If section or key isn't found,
// default value will be returned
func GetBool(section string, key string, def bool) bool {
    if value, err := get(section, key); err == nil {
        var tmp, _ = strconv.ParseBool(value)
        return tmp
    }

    return def
}

// Attempts to find and return given value from section and it's key,
// if section and key wasn't found, an error will be returned
func get(section string, key string) (string, error) {
    if sections[section] == nil {
        return "", errors.New("Cannot find section: " + section)
    }

    var data = sections[section]

    if data[key] == "" {
        return "", errors.New("Cannot find " + section + "::" + key)
    }

    log.Debugf("%s::%s=%s", section, key, data[key])
    return data[key], nil
}