package def

import (
	"github.com/ubis/Freya/cmd/loginserver/rsa"
	"github.com/ubis/Freya/share/models/server"
)

type Settings struct {
	server.Settings
	RSA rsa.RSA
}
