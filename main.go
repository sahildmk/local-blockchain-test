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

	GenesisChain := newChain("Main Chain")
	fmt.Println(GenesisChain.chain[0].PrevHash)
	fmt.Println(GenesisChain.chain[0].blockHash())

	// newT := Transaction{From:"Sahil", To:"John", Value: 10000}
	// newB := createBlock("0", newT)

	// newH := blockHash(newB)
	// fmt.Println(newH)
}