package models

import "share/encryption"

type Settings struct {
    ListenPort  int
    MaxUsers    int
    XorKeyTable encryption.XorKeyTable
}