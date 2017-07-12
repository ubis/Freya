package def

import (
    "loginserver/rsa"
    "share/encryption"
)

type Settings struct {
    XorKeyTable encryption.XorKeyTable
    RSA         rsa.RSA
}