package main

import (
	"fmt"
	"hexblox/internal/blockchain"
)

func main() {
	bchain := blockchain.NewBlockchain[string]()
	fmt.Println(bchain)
	a := "New block"
	bchain.AddBlock([]*string{&a})
	fmt.Println("-----------------")
	fmt.Println(bchain)
	fmt.Println(blockchain.IsValidChain(bchain.Chain()))
	block := blockchain.MineBlock(bchain.Chain()[1], []*string{&a})
	bchain.Chain()[1] = block
	fmt.Println(blockchain.IsValidChain(bchain.Chain()))
}
