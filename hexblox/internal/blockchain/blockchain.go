package blockchain

import "fmt"

type Blockchain[T any] struct {
	chain []*Block[T]
}

func NewBlockchain[T any]() *Blockchain[T] {
	return &Blockchain[T]{
		chain: []*Block[T]{Genesis[T]()},
	}
}

func (blockchain *Blockchain[T]) Chain() []*Block[T] {
	return blockchain.chain
}

func (blockchain *Blockchain[T]) AddBlock(data []*T) *Block[T] {
	block := MineBlock(blockchain.chain[len(blockchain.chain)-1], data)
	blockchain.chain = append(blockchain.chain, block)
	return block
}

func (blockchain *Blockchain[T]) String() string {
	chainString := "---Chain \n"
	for _, block := range blockchain.chain {
		chainString = fmt.Sprint(chainString, "   ", block)
	}
	return chainString
}

func IsValidChain[T any](chain []*Block[T]) bool {
	if chain[0].String() != Genesis[T]().String() {
		return false
	}

	for i := 1; i < len(chain); i++ {
		block := chain[i]
		lastBlock := chain[i-1]

		if block.LastHash != lastBlock.Hash {
			return false
		}
		if block.Hash != BlockHash(block) {
			return false
		}
	}

	return true
}

func (blockchain *Blockchain[T]) ReplaceChain(newChain []*Block[T]) {
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
