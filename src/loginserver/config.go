package main

import (
	"strconv"
	"share/lib/goini"
	"share/directory"
)

// Default values
const (
	C_Port 	 	= 38101
	C_MaxUsers 	= 100
)

// Configuration struct
type Config struct {
	file 		*goini.INI
	Port		int
	MaxUsers	int
}

// Attempts to read server configuration file
func (c *Config) Read() {
	log.Info("Reading configuration...")

	var location = directory.Root() + "/cfg/loginserver.ini"
	var ini 	 = goini.New()

	// parse configuration file...
	if err := ini.ParseFile(location); err != nil {
		log.Fatalf("Couldn't read configuration file %s. %s", location, err.Error())
		return
	}

	c.file = ini

	// read values from configuration...
	c.Port, _     = c.getInt("network", "port", C_Port)
	c.MaxUsers, _ = c.getInt("server", "max_users", C_MaxUsers)
}

/*
	Gets value from configuration file, if section or key isn't found,
	default value will be returned
	@param 	section 	conf section
	@param	key			conf key
	@param	def			default value
	@return	string value, either from conf or default and error, if any
 */
func (c *Config) getString(section string, key string, def string) (string, error) {
	if value, err := c.file.SectionGet(section, key); err != true {
		return value, nil
	}

	return def, nil
}

/*
	Gets value from configuration file, if section or key isn't found,
	default value will be returned
	@param 	section 	conf section
	@param	key			conf key
	@param	def			default value
	@return	int value, either from conf or default and error, if any
 */
func (c *Config) getInt(section string, key string, def int) (int, error) {
	if value, err := c.file.SectionGet(section, key); err != true {
		return strconv.Atoi(value)
	}

	return def, nil
}

/*
	Gets value from configuration file, if section or key isn't found,
	default value will be returned
	@param 	section 	conf section
	@param	key			conf key
	@param	def			default value
	@return	bool value, either from conf or default and error, if any
 */
func (c *Config) getBool(section string, key string, def bool) (bool, error) {
	if value, err := c.file.SectionGet(section, key); err != true {
		return strconv.ParseBool(value)
	}

	return def, nil
}