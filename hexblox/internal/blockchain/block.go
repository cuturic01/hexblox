package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hexblox/internal/config"
	"hexblox/internal/wallet"
	"strings"
	"time"
)

type Block struct {
	Timestamp  int64
	LastHash   string
	Hash       string
	Nonce      int
	Difficulty int
	Data       []*wallet.Transaction
}

func Genesis() *Block {
	return &Block{
		Timestamp:  0,
		LastHash:   "----------------",
		Hash:       "f1r57-h45h",
		Nonce:      0,
		Difficulty: config.Difficulty,
		Data:       []*wallet.Transaction{},
	}
}

func (block *Block) String() string {
	return fmt.Sprint(
		"-Block \n",
		"      Timestamp:  ", block.Timestamp, "\n",
		"      LastHash:   ", block.LastHash, "\n",
		"      Hash:       ", block.Hash, "\n",
		"      Nonce:      ", block.Nonce, "\n",
		"      Difficulty: ", block.Difficulty, "\n",
		"      Data:       ", block.Data, "\n",
	)
}

func MineBlock(lastBlock *Block, data []*wallet.Transaction) *Block {
	fmt.Println("Mining block...")

	var timestamp int64
	var hash string
	var difficulty int
	nonce := 0
	lastHash := lastBlock.Hash

	startTime := time.Now()
	for {
		timestamp = time.Now().UnixMilli()
		difficulty = adjustDifficulty(lastBlock, timestamp)

		var dataString string
		for _, transaction := range data {
			dataString += transaction.String()
		}
		hash = GenerateHash(timestamp, lastHash, dataString, nonce, difficulty)

		if hash[:difficulty] == strings.Repeat("0", difficulty) {
			endTime := time.Now()
			elapsedTime := endTime.Sub(startTime)
			fmt.Printf("Block successfully mined in: %fs.\n", elapsedTime.Seconds())
			return &Block{
				Timestamp:  timestamp,
				LastHash:   lastHash,
				Hash:       hash,
				Nonce:      nonce,
				Difficulty: difficulty,
				Data:       data,
			}
		}

		nonce++
	}
}

func GenerateHash(timestamp int64, lastHash string, data string, nonce int, difficulty int) string {
	s := fmt.Sprint(timestamp, lastHash, data, nonce, difficulty)
	hash := sha256.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

func HashBlock(block *Block) string {
	var dataString string
	for _, transaction := range block.Data {
		dataString += transaction.String()
	}
	return GenerateHash(block.Timestamp, block.LastHash, dataString, block.Nonce, block.Difficulty)
}

func adjustDifficulty(lastBlock *Block, currentTimestamp int64) int {
	if lastBlock.Timestamp+config.MineRate > currentTimestamp {
		return lastBlock.Difficulty + 1
	} else {
		return lastBlock.Difficulty - 1
	}
}
