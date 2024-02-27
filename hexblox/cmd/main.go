package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hexblox/internal/api/routes"
	"hexblox/internal/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain[string]()
	fmt.Println(bc)
	a := "New block"
	bc.AddBlock([]*string{&a})
	fmt.Println("-----------------")
	fmt.Println(bc)
	fmt.Println(blockchain.IsValidChain(bc.Chain()))

	r := gin.Default()
	routes.SetBlockchainRoutes(r, bc)
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
