package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type Account string

type Transaction struct {
	From 	Account 
	To 		Account
	Value 	uint
}

type Wallet struct {
	privateKey *rsa.PrivateKey
	publicKey rsa.PublicKey
}

func NewWallet() (Wallet, error) {
	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	checkError(err)

	publicKey := key.PublicKey

	return Wallet{privateKey: key, publicKey: publicKey}, nil
}

type Block struct {
	PrevHash string
	T Transaction
	TimeStamp string
}

func NewBlock(prevHash string, t Transaction) (Block) {
	return Block{PrevHash: prevHash, T: t, TimeStamp: time.Now().String()}
}

type State struct {
	Balances map[Account]uint
	dbFile *os.File
}

type Chain struct {
	chain []Block
	chainID string
}

func (c Chain) getLastBlock() (Block) {
	return c.chain[len(c.chain) - 1]
}

// sPK => Senders Public Key
// sig => Signature
func (c Chain) addBlock(t Transaction, sPK string, sig string) {

}

func newChain(chainID string) (Chain) {
	token := []byte(time.Now().String())
	hS := sha256.Sum256(token)
	hash := fmt.Sprintf("%x", hS[:])
	newB := NewBlock(hash, Transaction{"Genesis", "Satoshi", 1000})
	var blocks []Block
	blocks = append(blocks, newB)
	return Chain{chain: blocks, chainID: chainID}
}

func (b Block) blockHash() (string) {
	s := fmt.Sprintf("%v", b)
	hS := sha256.Sum256([]byte(s))
	hash := fmt.Sprintf("%x", hS[:])
	return hash
}

func transactionToString(t Transaction) (string, error) {

	return "hi", nil
}

func loadToMap(path string) (map[string]interface{}, error) {
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Printf("Found %s", path)

	bVal, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
    json.Unmarshal([]byte(bVal), &result)

	return result, nil
}

func stateFromDisk() (*State, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	dbFilePath := filepath.Join(cwd, "data", "genesis.json")
	// Load genesis database
	
	if err != nil {
		return nil, err
	}

	var genFile map[string]interface{}
	genFile, err = loadToMap(dbFilePath)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(genFile["gen_date"])

	return nil, nil
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func main() {
	// cwd, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// dbFilePath := filepath.Join(cwd, "data", "genesis.json")
	// // Load genesis database
	
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// var genFile map[string]interface{}
	// genFile, err = loadToMap(dbFilePath)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// GenesisChain := newChain("Main Chain")
	// fmt.Println(GenesisChain.chain[0].PrevHash)

	w, err := NewWallet()
	if err != nil {
		fmt.Println(err)
	}

	testVar := "Hey hows it going"

	enc, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &w.publicKey, []byte(testVar), nil)
	dec, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, w.privateKey, enc, nil)

	fmt.Println(string(dec))

	// newT := Transaction{From:"Sahil", To:"John", Value: 10000}
	// newB := createBlock("0", newT)

	// newH := blockHash(newB)
	// fmt.Println(newH)
}