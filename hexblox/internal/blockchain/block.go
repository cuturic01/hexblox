package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block[T any] struct {
	timestamp int64
	lastHash  string
	hash      string
	data      []*T
}

func Genesis[T any]() *Block[T] {
	return &Block[T]{
		timestamp: 0,
		lastHash:  "----------------",
		hash:      "f1r57-h45h",
		data:      []*T{},
	}
}

func NewBlock[T any](timestamp int64, lastHash string, hash string, data []*T) *Block[T] {
	return &Block[T]{
		timestamp: timestamp,
		lastHash:  lastHash,
		hash:      hash,
		data:      data,
	}
}

func (block *Block[T]) Data() []*T {
	return block.data
}

func (block *Block[T]) String() string {
	return fmt.Sprint(
		"-Block \n",
		"      timestamp: ", block.timestamp, "\n",
		"      lastHash:  ", block.lastHash[0:10], "...\n",
		"      hash:      ", block.hash[0:10], "...\n",
		"      data:      ", block.data, "\n",
	)
}

func MineBlock[T any](lastBlock *Block[T], data []*T) *Block[T] {
	timestamp := time.Now().UnixMilli()
	lastHash := lastBlock.hash
	hash := GenerateHash(timestamp, lastHash, data)

	return &Block[T]{
		timestamp: timestamp,
		lastHash:  lastHash,
		hash:      hash,
		data:      data,
	}
}

func GenerateHash[T any](timestamp int64, lastHash string, data []T) string {
	s := fmt.Sprint(timestamp, lastHash, data)
	hash := sha256.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

func BlockHash[T any](block *Block[T]) string {
	return GenerateHash(block.timestamp, block.lastHash, block.data)
}
