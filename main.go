package main

import (
	"bufio"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	mrand "math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var mainChain = newChain("Genesis Chain")
const TXPBLOCK = 5
var aliasToPKMap = make(map[string]rsa.PublicKey)
var PKToAliasMap = make(map[rsa.PublicKey]string)
var aliasToWallet = make(map[string]Wallet)
var users = []string{}

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
			fmt.Printf("%s		%s		%.2f\n", PKToAliasMap[t.To], PKToAliasMap[t.From], t.Value)
		}
		fmt.Println()
	}
}

func printBal(user string) () {
	if w, exists := aliasToWallet[user]; exists {
		fmt.Printf("%s's Balance: %.2f\n", user, w.balance)
	} else {
		fmt.Println("User does not exist!")
	}
}

func printUsers() () {
	fmt.Println("Current Users")
	for i := 0; i < len(users); i++ {
		fmt.Printf("- %s\n", users[i])
	}
}

func printBlock() () {
	b := mainChain.getLastBlock()
	fmt.Printf("> BLOCK %d\n", len(mainChain.chain))
	fmt.Printf("Length of Block: %d\n", b.getTLen())
	fmt.Println("TO		FROM 		VALUE")
	for i := 0; i < len(b.transList); i++ {
		t := b.transList[i]
		fmt.Printf("%s		%s		%.2f\n", PKToAliasMap[t.To], PKToAliasMap[t.From], t.Value)
	}
}

func userExists(user string) (bool) {
	for i := 0; i < len(users); i++ {
		if users[i] == user {
			return true
		}
	}
	return false
}

// =============== Transaction Functions =============== //
type Transaction struct {
	From 	rsa.PublicKey 
	To 		rsa.PublicKey
	Value 	float32
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
	balance float32
	privateKey *rsa.PrivateKey
	publicKey rsa.PublicKey
}

func NewWallet(name string) (error) {

	if _, exists := aliasToPKMap[name]; exists {
		return errors.New("user already exists")
	}

	users = append(users, name)

	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	checkError(err)

	publicKey := key.PublicKey

	aliasToPKMap[name] = publicKey
	PKToAliasMap[publicKey] = name

	w := Wallet{privateKey: key, publicKey: publicKey, balance: 1000}
	aliasToWallet[name] = w

	return nil
}

func (w Wallet) sendMoney(amount float32, target string) (error) {
	// fmt.Println("> Sending Money")

	if (amount > w.balance) {
		return errors.New("insufficient funds")
	}

	var targetPK rsa.PublicKey

	if pk, exists := aliasToPKMap[target]; exists {
		targetPK = pk
	} else {
		return errors.New("> Target does not have a wallet")
	}

	newT := Transaction{From: w.publicKey, To: targetPK, Value: amount}
	hashed := sha256.Sum256([]byte(newT.transToString()))
	sig, err := rsa.SignPKCS1v15(rand.Reader, w.privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Println("========== Error from signing: ", err)
		return errors.New("signing error")
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

		from := aliasToWallet[PKToAliasMap[t.From]]
		to := aliasToWallet[PKToAliasMap[t.To]]

		from.balance -= t.Value
		to.balance += t.Value

		aliasToWallet[PKToAliasMap[t.From]] = from
		aliasToWallet[PKToAliasMap[t.To]] = to

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

	gErr := NewWallet("Genesis")
	sErr := NewWallet("Satoshi")

	genesis := aliasToWallet["Genesis"]
	satoshi := aliasToWallet["Satoshi"]

	if gErr != nil {
		panic(gErr)
	} else if sErr != nil {
		panic(sErr)
	}

	newB := NewBlock(hash)
	genesis.balance = 1000000
	aliasToWallet["Genesis"] = genesis
	newB, err := newB.addTx(Transaction{genesis.publicKey, satoshi.publicKey, 1000})
	if err != nil {
		panic(err)
	}
	var blocks []Block
	blocks = append(blocks, newB)
	return Chain{chain: blocks, chainID: chainID}
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
			newB.addTx(t)
			c.chain = append(c.chain, newB)
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func main() {
	mrand.Seed(time.Now().UnixNano())

	fmt.Println(
		"List of Commands\n",
		"h - Help\n",
		"q - Quit\n",
		"new [username] - Creates a new wallet\n",
		"send [fromUser] [toUser] [value] - Send money to a user\n",
		"bal [username] - Shows a users balance\n",
		"users - Prints a list of all current users\n",
		"block - Prints the transactions of the current block\n",
		"chain - Prints the entire chain",
	)

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		splitInput := strings.Split(input, " ")	

		switch splitInput[0] {
		case "q":
			fmt.Print("Data is not persistent, are you sure you want to exit? (Y/N) ")
			scanner.Scan()
			input := scanner.Text()
			if input == "Y" || input == "y" {
				fmt.Println("Goodbye!")
				os.Exit(0)
			}
		case "h":
			fmt.Println(
				"List of Commands\n",
				"h - Help\n",
				"q - Quit\n",
				"new [username] - Creates a new wallet\n",
				"send [fromUser] [toUser] [value] - Send money to a user\n",
				"bal [username] - Shows a users balance\n",
				"users - Prints a list of all current users\n",
				"block - Prints the transactions of the current block\n",
				"chain - Prints the entire chain",
			)
		case "new":
			name := splitInput[1]
			err := NewWallet(name)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("New wallet created for %s!\n", name)
		case "chain":
			printChain()
		case "bal":
			printBal(splitInput[1])
		case "users":
			printUsers()
		case "block":
			printBlock()
		case "send":
			fromUser := splitInput[1]
			toUser := splitInput[2]
			value, err := strconv.ParseFloat(splitInput[3], 32)
			if (err != nil) {
				panic(err)
			}
			if userExists(fromUser) {
				if userExists(toUser) {
					fromWallet := aliasToWallet[fromUser]
					err := fromWallet.sendMoney(float32(value), toUser)
					if (err == nil) {
						fmt.Printf("Sent %.2f from %s to %s\n", value, fromUser, toUser)
					} else {
						fmt.Println(err)
					}
				} else {
					fmt.Printf("User %s does not exist\n", toUser)
				}
			} else {
				fmt.Printf("User %s does not exist\n", fromUser)
			}
		default:
			fmt.Println("Unknown command. Type 'h' for a list of commands.")
		}

	}
}