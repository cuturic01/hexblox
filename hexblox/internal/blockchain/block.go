package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	Timestamp int64
	LastHash  string
	Hash      string
	Data      []*string
	Nonce     int
}

func Genesis() *Block {
	return &Block{
		Timestamp: 0,
		LastHash:  "----------------",
		Hash:      "f1r57-h45h",
		Data:      []*string{},
		Nonce:     0,
	}
}

func (block *Block) String() string {
	return fmt.Sprint(
		"-Block \n",
		"      Timestamp: ", block.Timestamp, "\n",
		"      LastHash:  ", block.LastHash[0:10], "...\n",
		"      Hash:      ", block.Hash[0:10], "...\n",
		"      Data:      ", block.Data, "\n",
		"      Nonce:     ", block.Nonce, "\n",
	)
}

func MineBlock(lastBlock *Block, data []*string, difficulty int) *Block {
	fmt.Println("Mining block...")

	var timestamp int64
	var hash string
	nonce := 0
	lastHash := lastBlock.Hash

	for {
		timestamp = time.Now().UnixMilli()
		hash = GenerateHash(timestamp, lastHash, data, nonce)

		if hash[:difficulty] == strings.Repeat("0", difficulty) {
			break
		}

		nonce++
	}

	fmt.Println("Block successfully mined.")

	return &Block{
		Timestamp: timestamp,
		LastHash:  lastHash,
		Hash:      hash,
		Data:      data,
		Nonce:     nonce,
	}
}

func GenerateHash(timestamp int64, lastHash string, data []*string, nonce int) string {
	var stringData string

	for _, strPtr := range data {
		stringData += *strPtr
	}

	s := fmt.Sprint(timestamp, lastHash, stringData, nonce)
	hash := sha256.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

func HashBlock(block *Block) string {
	return GenerateHash(block.Timestamp, block.LastHash, block.Data, block.Nonce)
}
