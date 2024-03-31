package domain

import (
	"fmt"
)

type Blockchain struct {
	chain []*Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		chain: []*Block{Genesis()},
	}
}

func (blockchain *Blockchain) Chain() []*Block {
	return blockchain.chain
}

func (blockchain *Blockchain) AddBlock(data []*Transaction) *Block {
	block := MineBlock(blockchain.chain[len(blockchain.chain)-1], data)
	blockchain.chain = append(blockchain.chain, block)
	return block
}

func (blockchain *Blockchain) String() string {
	chainString := "---Chain \n"
	for _, block := range blockchain.chain {
		chainString = fmt.Sprint(chainString, "   ", block)
	}
	return chainString
}

func IsValidChain(chain []*Block) bool {
	if chain[0].String() != Genesis().String() {
		return false
	}

	for i := 1; i < len(chain); i++ {
		block := chain[i]
		lastBlock := chain[i-1]

		if block.LastHash != lastBlock.Hash {
			return false
		}
		if block.Hash != HashBlock(block) {
			return false
		}
	}

	return true
}

func (blockchain *Blockchain) ReplaceChain(newChain []*Block) {
	if len(blockchain.chain) >= len(newChain) {
		fmt.Println("Received chain is not longer then the current.")
		return
	}
	if !IsValidChain(newChain) {
		fmt.Println("Chain is not valid!")
		return
	}

	fmt.Println("Replacing current chain with the new chain.")
	blockchain.chain = newChain
}
