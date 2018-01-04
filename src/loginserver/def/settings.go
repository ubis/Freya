package def

import (
	"loginserver/rsa"
	"share/models/server"
)

type Settings struct {
	server.Settings
	RSA rsa.RSA
}
