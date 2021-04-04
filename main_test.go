package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestSigning(t *testing.T) {
	fmt.Println("> Creating Wallet 1")
	w1, err := NewWallet()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("> Wallet 1 Created!")
	
	fmt.Println("> Creating Wallet 2")
	w2, err := NewWallet()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("> Wallet 2 Created!")

	newT := Transaction{From: w1.publicKey, To: w2.publicKey, Value: 1000}
	hashed := sha256.Sum256([]byte(newT.transToString()))
	
	fmt.Println("> Signing test transaction")
	sig, err := rsa.SignPKCS1v15(rand.Reader, w1.privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		t.Errorf( "Error from signing: %s\n", err)
	}

	fmt.Println("> Verifying test transaction")
	verifyErr := rsa.VerifyPKCS1v15(&w1.publicKey, crypto.SHA256, hashed[:], sig)
	if verifyErr != nil {
		t.Errorf("Error from verification: %s\n", err)
	} else {
		fmt.Println("> Transaction has been verified")
	}
}