package Chain

import (
	"crypto/sha256"
	"fmt"
	"local-blockchain-test/Block"
	"time"
)

type Chain struct {
	chain []Block
	chainID string
}

func (c Chain) getLastBlock() (Block) {
	return c.chain[len(c.chain) - 1]
}

// sPK => Senders Public Key
// sig => Signature
func (c Chain) addBlock(t Transaction, sPK string, sig string)

func newChain(chainID string) (Chain) {
	token := []byte(time.Now().String())
	hS := sha256.Sum256(token)
	hash := fmt.Sprintf("%x", hS[:])
	newB := Block.NewBlock(hash, Transaction{"Genesis", "Satoshi", 1000})
	var blocks []Block
	blocks = append(blocks, newB)
	return Chain{chain: blocks, chainID: chainID}
}