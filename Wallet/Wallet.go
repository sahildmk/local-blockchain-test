package Wallet

import (
	"crypto/rand"
	"crypto/rsa"
)

type Wallet struct {
	KeyPair *rsa.PrivateKey
}

func NewWallet() *Wallet {
	w := new(Wallet)
	keyPair, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		w.KeyPair = keyPair
	}
	return w
}