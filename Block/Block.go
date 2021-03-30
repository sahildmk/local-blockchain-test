package Block

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Block struct {
	PrevHash string
	T Transaction
	TimeStamp string
}

func NewBlock(prevHash string, t Transaction) *Block {
	b := new(Block)
	b.PrevHash = prevHash
	b.T = t
	b.TimeStamp = time.Now().String() 
	return b
}

func (b Block) blockHash() (string) {
	s := fmt.Sprintf("%v", b)
	hS := sha256.Sum256([]byte(s))
	hash := fmt.Sprintf("%x", hS[:])
	return hash
}