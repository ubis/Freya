package rsa

import (
    "crypto/rsa"
    "crypto/rand"
    "crypto/sha1"
    "encoding/asn1"
    "share/logger"
)

const RSA_KEY_LENGTH     = 2048
const RSA_LOGIN_LENGTH   = RSA_KEY_LENGTH / 8
const RSA_PUB_KEY_LENGTH = RSA_LOGIN_LENGTH + 14

var log = logger.Instance()

type RSA struct {
    privateKey *rsa.PrivateKey

    PublicKey  [RSA_PUB_KEY_LENGTH]uint8
}

// Initializes RSA which generates keypair
func (r *RSA) Init() {
    log.Infof("Generating %d bit RSA key...", RSA_KEY_LENGTH)

    // generate RSA key
    var err error
    r.privateKey, err = rsa.GenerateKey(rand.Reader, RSA_KEY_LENGTH)
    if err != nil {
        log.Error("Error generating RSA key: " + err.Error())
        return
    }

    // encode key to ASN.1 PublicKey Type
    var key, err2 = asn1.Marshal(r.privateKey.PublicKey)
    if err2 != nil {
        log.Error("Error encoding Public RSA key: " + err2.Error())
        return
    }

    // move public key to array
    copy(r.PublicKey[:], key)
}

/*
    Attempts to decrypt RSA data, which is `RSA_LOGIN_LENGTH` length
    @param  data    data array to be decrypted
    @return decrypted data or error on fail
 */
func (r *RSA) Decrypt(data []uint8) ([]byte, error) {
    var dec, err = rsa.DecryptOAEP(sha1.New(), rand.Reader, r.privateKey, data, nil)
    if err != nil {
        return nil, err
    }

    return dec, nil
}