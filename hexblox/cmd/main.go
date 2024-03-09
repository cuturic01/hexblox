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
	transaction2 := wallet.NewTransaction(wallet2, wallet1.PublicKey, 60)
	transactionPool.AddTransaction(transaction)
	transactionPool.AddTransaction(transaction2)
	fmt.Println("Before:", transactionPool)
	transaction.Update(wallet1, wallet2.PublicKey, 100)
	transactionPool.AddTransaction(transaction)
	fmt.Println("After:", transactionPool)

	//httpPort := flag.String("http", "", "HTTP server port")
	//hostPort := flag.String("host", "", "WebSocket server port")
	//flag.Parse()
	//
	//p2p.Run(*httpPort, *hostPort)
	//
	//select {}
}
