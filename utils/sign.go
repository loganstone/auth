package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"log"

	"gopkg.in/square/go-jose.v2"
)

var privateKey rsa.PrivateKey
var signer jose.Signer

func init() {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Panicln(err)
	}
	privateKey = *key

	object, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.PS512,
			Key:       key,
		}, nil)
	if err != nil {
		log.Panicln(err)
	}
	signer = object
}

// Sign .
func Sign(payload []byte) (string, error) {
	object, err := signer.Sign(payload)
	if err != nil {
		return "", err
	}
	return object.CompactSerialize()
}

// Load .
func Load(serialized string) ([]byte, error) {
	object, err := jose.ParseSigned(serialized)
	if err != nil {
		return nil, err
	}
	return object.Verify(&privateKey.PublicKey)
}
