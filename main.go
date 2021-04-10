package main

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	mrand "math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var mainChain = newChain("Genesis Chain")
const TXPBLOCK = 5
var aliasToPKMap = make(map[string]rsa.PublicKey)
var PKToAliasMap = make(map[rsa.PublicKey]string)

// =============== Helper Functions =============== //
func printChain() () {
	fmt.Printf("Chain Length: %d\n", len(mainChain.chain))

	for i := 0; i < len(mainChain.chain); i++ {
		var b Block = mainChain.chain[i]
		fmt.Printf("> BLOCK %d\n", i)
		fmt.Printf("Length of Block: %d\n", b.getTLen())
		fmt.Println("Transactions in Block")
		fmt.Println("TO		FROM 		VALUE")
		for j := 0; j < len(b.transList); j++ {
			var t Transaction = b.transList[j]
			fmt.Printf("%s		%s		%d\n", PKToAliasMap[t.To], PKToAliasMap[t.From], t.Value)
		}
		fmt.Println()
	}
}


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

func NewWallet(name string) (Wallet, error) {

	if _, exists := aliasToPKMap[name]; exists {
		return Wallet{}, errors.New("User already exists")
	}

	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	checkError(err)

	publicKey := key.PublicKey

	aliasToPKMap[name] = publicKey
	PKToAliasMap[publicKey] = name

	return Wallet{privateKey: key, publicKey: publicKey}, nil
}

func (w Wallet) sendMoney(amount uint, target string) (error) {
	// fmt.Println("> Sending Money")
	var targetPK rsa.PublicKey

	if pk, exists := aliasToPKMap[target]; exists {
		targetPK = pk
	} else {
		return errors.New("> Target does not have a wallet!")
	}

	newT := Transaction{From: w.publicKey, To: targetPK, Value: amount}
	hashed := sha256.Sum256([]byte(newT.transToString()))
	sig, err := rsa.SignPKCS1v15(rand.Reader, w.privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Println("========== Error from signing: ", err)
		return errors.New("> Signing error")
	}

	mainChain.newTrans(newT, w.publicKey, sig)

	return nil
}

// =============== Block Functions =============== //
type Block struct {
	nonce int
	PrevHash string
	transList []Transaction
	TimeStamp string
}

func NewBlock(prevHash string) (Block) {
	// fmt.Println("> Creating a new block")
	n := int(mrand.Float64()*999999999)
	var tList []Transaction
	return Block{PrevHash: prevHash, transList: tList, TimeStamp: time.Now().String(), nonce: n}
}

func (b Block) blockHash() (string) {
	s := fmt.Sprintf("%v", b)
	hS := sha256.Sum256([]byte(s))
	hash := fmt.Sprintf("%x", hS[:])
	return hash
}

func (b Block) addTx(t Transaction) (Block, error) {
	if b.getTLen() < TXPBLOCK {
		// fmt.Println("> Adding transaction", len(b.transList))
		b.transList = append(b.transList, t)
		return b, nil
	}
	return b, errors.New("Block full")
}

func (b Block) getTLen() (int) {
	return len(b.transList)
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

	genesis, gErr := NewWallet("Genesis")
	satoshi, sErr := NewWallet("Satoshi")

	if gErr != nil {
		panic(gErr)
	} else if sErr != nil {
		panic(sErr)
	}

	newB := NewBlock(hash)
	newB, err := newB.addTx(Transaction{genesis.publicKey, satoshi.publicKey, 1000})
	if err != nil {
		panic(err)
	}
	var blocks []Block
	blocks = append(blocks, newB)
	return Chain{chain: blocks, chainID: chainID}
}
type State struct {
	Balances map[string]uint
	transactions []Transaction
	dbFile *os.File
}

func (c Chain) getLastBlock() (Block) {
	return c.chain[len(c.chain) - 1]
}

// sPK => Senders Public Key
// sig => Signature
func (c *Chain) newTrans(t Transaction, sPK rsa.PublicKey, sig []byte) () {
	// fmt.Println("> New Transaction")
	hashed := sha256.Sum256([]byte(t.transToString()))
	err := rsa.VerifyPKCS1v15(&sPK, crypto.SHA256, hashed[:], sig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from verification: %s\n", err)
    	return
	} else {
		// fmt.Println("> Transaction has been verified")
		c.chain[len(c.chain) - 1], err = c.getLastBlock().addTx(t)
		if err != nil {
			newB := NewBlock(c.getLastBlock().blockHash())
			newB, err = newB.addTx(t)
			c.chain = append(c.chain, newB)
		}
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

	for {
		var input string
		fmt.Printf("> ")
		fmt.Scanln(&input)		
		
		switch input {
		case "q":
			fmt.Println("Goodbye!")
			os.Exit(0)
		default:
			fmt.Println("Unknown command. Type 'h' for a list of commands.")
		}

	}

	w1, err := NewWallet("sahil")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println("> Wallet 1 Created!")
	
	// fmt.Println("> Creating Wallet 2")
	w2, err := NewWallet("john")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println("> Wallet 2 Created!")

	for i := 0; i < 6; i++ {
		w1.sendMoney(1000, "john")
	}

	for i := 0; i < 6; i++ {
		w2.sendMoney(1000, "sahil")
	}

	
	printChain()
	

	// newB := createBlock("0", newT)

	// newH := blockHash(newB)
	// fmt.Println(newH)
	// fmt.Println(fmt.Sprintf("%.9f", mrand.Float64()))
}