package utils

import (
	"crypto/rand"
	"crypto/rsa"

	"gopkg.in/square/go-jose.v2"
)

var privateKey rsa.PrivateKey
var signer jose.Signer

func init() {
	// Generate a public/private key pair to use for this example.
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	privateKey = *key

	// Instantiate a signer using RSASSA-PSS (SHA512) with the given private key.
	object, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.PS512,
			Key:       key,
		}, nil)
	if err != nil {
		panic(err)
	}
	signer = object
}

// Sign .
func Sign(payload []byte) (string, error) {
	object, err := signer.Sign(payload)
	if err != nil {
		panic(err)
	}
	return object.CompactSerialize()
}

// Load .
func Load(serialized string) []byte {
	object, err := jose.ParseSigned(serialized)
	if err != nil {
		panic(err)
	}
	output, err := object.Verify(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}
	return output
}
