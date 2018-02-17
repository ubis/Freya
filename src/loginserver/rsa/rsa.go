package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/asn1"
	"share/log"
)

const (
	keyLength       = 2048
	LoginLength     = keyLength / 8
	publicKeyLength = LoginLength + 14
)

// RSA Encryption structure
type Encryption struct {
	privateKey *rsa.PrivateKey
	PublicKey  [publicKeyLength]byte
}

// Initializes RSA which generates keypair
func (r *Encryption) Init() {
	log.Infof("Generating %d bit RSA key...", keyLength)

	// generate RSA key
	pKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		log.Fatal("Error generating RSA key:" + err.Error())
	}

	// store private key
	r.privateKey = pKey

	// encode key to ASN.1 PublicKey Type
	key, err := asn1.Marshal(r.privateKey.PublicKey)
	if err != nil {
		log.Fatal("Error encoding Public RSA key:" + err.Error())
	}

	// move public key to array
	copy(r.PublicKey[:], key)
}

// Attempts to decrypt RSA data, which is `LoginLength` length
func (r *Encryption) Decrypt(data []byte) ([]byte, error) {
	hash := sha1.New()
	reader := rand.Reader
	dec, err := rsa.DecryptOAEP(hash, reader, r.privateKey, data, nil)
	if err != nil {
		return nil, err
	}

	return dec, nil
}
