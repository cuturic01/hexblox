package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Block[T any] struct {
	Timestamp int64
	LastHash  string
	Hash      string
	Data      []*T
}

func Genesis[T any]() *Block[T] {
	return &Block[T]{
		Timestamp: 0,
		LastHash:  "----------------",
		Hash:      "f1r57-h45h",
		Data:      []*T{},
	}
}

func (block *Block[T]) String() string {
	return fmt.Sprint(
		"-Block \n",
		"      Timestamp: ", block.Timestamp, "\n",
		"      LastHash:  ", block.LastHash[0:10], "...\n",
		"      Hash:      ", block.Hash[0:10], "...\n",
		"      Data:      ", block.Data, "\n",
	)
}

func MineBlock[T any](lastBlock *Block[T], data []*T) *Block[T] {
	timestamp := time.Now().UnixMilli()
	lastHash := lastBlock.Hash
	hash := GenerateHash(timestamp, lastHash, data)

	return &Block[T]{
		Timestamp: timestamp,
		LastHash:  lastHash,
		Hash:      hash,
		Data:      data,
	}
}

func GenerateHash[T any](timestamp int64, lastHash string, data []T) string {
	s := fmt.Sprint(timestamp, lastHash, data)
	hash := sha256.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

func BlockHash[T any](block *Block[T]) string {
	return GenerateHash(block.Timestamp, block.LastHash, block.Data)
}
