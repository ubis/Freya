package def

import (
    "share/models/server"
    "loginserver/rsa"
)

type Settings struct {
    server.Settings
    RSA         rsa.RSA
}