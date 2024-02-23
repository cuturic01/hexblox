package main

import (
	"fmt"
	"hexblox/internal/blockchain"
	"time"
)

func main() {
	genesis := blockchain.Genesis()
	block := blockchain.New(time.Now().UnixMilli(), "ax001pdsa210000", "ax012234500000", []string{"Data"})
	fmt.Println(genesis)
	fmt.Println(block)
	newBlock := blockchain.MineBlock(block, []string{"New block"})
	fmt.Println(newBlock)
}
