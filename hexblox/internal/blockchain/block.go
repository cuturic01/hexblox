package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block struct {
	Timestamp int64
	LastHash  string
	Hash      string
	Data      []*string
}

func Genesis() *Block {
	return &Block{
		Timestamp: 0,
		LastHash:  "----------------",
		Hash:      "f1r57-h45h",
		Data:      []*string{},
	}
}

func (block *Block) String() string {
	return fmt.Sprint(
		"-Block \n",
		"      Timestamp: ", block.Timestamp, "\n",
		"      LastHash:  ", block.LastHash[0:10], "...\n",
		"      Hash:      ", block.Hash[0:10], "...\n",
		"      Data:      ", block.Data, "\n",
	)
}

func MineBlock(lastBlock *Block, data []*string) *Block {
	timestamp := time.Now().UnixMilli()
	lastHash := lastBlock.Hash
	hash := GenerateHash(timestamp, lastHash, data)

	return &Block{
		Timestamp: timestamp,
		LastHash:  lastHash,
		Hash:      hash,
		Data:      data,
	}
}

func GenerateHash(timestamp int64, lastHash string, data []*string) string {
	var stringData string

	for _, strPtr := range data {
		stringData += *strPtr
	}

	s := fmt.Sprint(timestamp, lastHash, stringData)
	hash := sha256.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

func HashBlock(block *Block) string {
	return GenerateHash(block.Timestamp, block.LastHash, block.Data)
}
