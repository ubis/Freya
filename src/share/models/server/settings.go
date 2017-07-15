package server

import "share/encryption"

type Settings struct {
    XorKeyTable  encryption.XorKeyTable
    CurrentUsers int
}