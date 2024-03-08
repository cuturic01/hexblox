package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hexblox/internal/config"
	"strings"
	"time"
)

type Block struct {
	Timestamp  int64
	LastHash   string
	Hash       string
	Nonce      int
	Difficulty int
	Data       []*string
}

func Genesis() *Block {
	return &Block{
		Timestamp:  0,
		LastHash:   "----------------",
		Hash:       "f1r57-h45h",
		Nonce:      0,
		Difficulty: config.Difficulty,
		Data:       []*string{},
	}
}

func (block *Block) String() string {
	return fmt.Sprint(
		"-Block \n",
		"      Timestamp:  ", block.Timestamp, "\n",
		"      LastHash:   ", block.LastHash[0:10], "...\n",
		"      Hash:       ", block.Hash[0:10], "...\n",
		"      Nonce:      ", block.Nonce, "\n",
		"      Difficulty: ", block.Difficulty, "\n",
		"      Data:       ", block.Data, "\n",
	)
}

func MineBlock(lastBlock *Block, data []*string) *Block {
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
		hash = GenerateHash(timestamp, lastHash, data, nonce, difficulty)

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

func GenerateHash(timestamp int64, lastHash string, data []*string, nonce int, difficulty int) string {
	var stringData string

	for _, strPtr := range data {
		stringData += *strPtr
	}

	s := fmt.Sprint(timestamp, lastHash, stringData, nonce, difficulty)
	hash := sha256.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

func HashBlock(block *Block) string {
	return GenerateHash(block.Timestamp, block.LastHash, block.Data, block.Nonce, block.Difficulty)
}

func adjustDifficulty(lastBlock *Block, currentTimestamp int64) int {
	if lastBlock.Timestamp+config.MineRate > currentTimestamp {
		return lastBlock.Difficulty + 1
	} else {
		return lastBlock.Difficulty - 1
	}
}
