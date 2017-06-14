package packet

import "share/encryption"

const (
    // get magic key from encryption package
    MAGIC_KEY               = encryption.MagicKey
)

const (
    CONNECT2SVR             = 101
    VERIFYLINKS             = 102
    AUTHACCOUNT             = 103
    SYSTEMMESSG             = 120
    SERVERSTATE             = 121
    CHECKVERSION            = 122
    URLTOCLIENT             = 128
    PUBLIC_KEY              = 2001
    PRE_SERVER_ENV_REQUEST  = 2002
)