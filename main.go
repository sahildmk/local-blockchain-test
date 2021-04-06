package main

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	mrand "math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Account string

var mainChain = newChain("Genesis Chain")

// =============== Helper Functions =============== //

// =============== Transaction Functions =============== //
type Transaction struct {
	From 	rsa.PublicKey 
	To 		rsa.PublicKey
	Value 	uint
}

func (t Transaction) transHash() (string) {
	s := fmt.Sprintf("%v", t)
	hS := sha256.Sum256([]byte(s))
	hash := fmt.Sprintf("%x", hS[:])
	return hash
}

func (t Transaction) transToString() (string) {
	out, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// =============== Wallet Functions =============== //
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

func (w Wallet) sendMoney(amount uint, targetPK rsa.PublicKey) () {
	newT := Transaction{From: w.publicKey, To: targetPK, Value: amount}
	hashed := sha256.Sum256([]byte(newT.transToString()))
	sig, err := rsa.SignPKCS1v15(rand.Reader, w.privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
		return
	}

	mainChain.addBlock(newT, w.publicKey, sig)

	return
}

// =============== Block Functions =============== //
type Block struct {
	nonce int
	PrevHash string
	T Transaction
	TimeStamp string
}

func NewBlock(prevHash string, t Transaction) (Block) {
	n := int(mrand.Float64()*999999999)
	return Block{PrevHash: prevHash, T: t, TimeStamp: time.Now().String(), nonce: n}
}

func (b Block) blockHash() (string) {
	s := fmt.Sprintf("%v", b)
	hS := sha256.Sum256([]byte(s))
	hash := fmt.Sprintf("%x", hS[:])
	return hash
}

// =============== Chain Functions =============== //
type Chain struct {
	chain []Block
	chainID string
}

func (c Chain) mine(nonce int) (int) {
	for sol := 0; true; sol++ {
		bs := []byte(strconv.Itoa(nonce + sol))
		hash := md5.Sum(bs)
		sHash := fmt.Sprintf("%x", hash)
		if (sHash[0:4] == "0000") {
			fmt.Println("Solved: ", sol)
			return sol
		}
	}
	return -1
}

func newChain(chainID string) (Chain) {
	token := []byte(time.Now().String())
	hS := sha256.Sum256(token)
	hash := fmt.Sprintf("%x", hS[:])

	genesis, gErr := NewWallet()
	satoshi, sErr := NewWallet()

	if gErr != nil {
		panic(gErr)
	} else if sErr != nil {
		panic(sErr)
	}

	newB := NewBlock(hash, Transaction{genesis.publicKey, satoshi.publicKey, 1000})
	var blocks []Block
	blocks = append(blocks, newB)
	return Chain{chain: blocks, chainID: chainID}
}
type State struct {
	Balances map[Account]uint
	dbFile *os.File
}

func (c Chain) getLastBlock() (Block) {
	return c.chain[len(c.chain) - 1]
}

// sPK => Senders Public Key
// sig => Signature
func (c Chain) addBlock(t Transaction, sPK rsa.PublicKey, sig []byte) {
	hashed := sha256.Sum256([]byte(t.transToString()))
	err := rsa.VerifyPKCS1v15(&sPK, crypto.SHA256, hashed[:], sig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from verification: %s\n", err)
    	return
	} else {
		fmt.Println("> Transaction has been verified")
		newB := NewBlock(c.getLastBlock().PrevHash, t)
		c.chain = append(c.chain, newB)
	}
}


// =============== Other Functions =============== //

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
	mrand.Seed(time.Now().UnixNano())
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

	// newB := createBlock("0", newT)

	// newH := blockHash(newB)
	// fmt.Println(newH)
	// fmt.Println(fmt.Sprintf("%.9f", mrand.Float64()))
}