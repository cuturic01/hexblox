package main

import (
	"fmt"
	"hexblox/internal/wallet"
)

func main() {

	transactionPool := wallet.NewTransactionPool()
	wallet1 := wallet.NewWallet()
	wallet2 := wallet.NewWallet()

	transaction := wallet.NewTransaction(wallet1, wallet2.PublicKey, 60)
	transactionPool.AddTransaction(transaction)
	fmt.Println("Before:", transactionPool)

	wallet1.CreateTransaction(wallet2.PublicKey, 60, transactionPool)
	wallet2.CreateTransaction(wallet1.PublicKey, 100, transactionPool)

	fmt.Println("After:", transactionPool)

	//httpPort := flag.String("http", "", "HTTP server port")
	//hostPort := flag.String("host", "", "WebSocket server port")
	//flag.Parse()
	//
	//p2p.Run(*httpPort, *hostPort)
	//
	//select {}
}
