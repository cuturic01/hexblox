package main

import (
	"fmt"
	"hexblox/internal/wallet"
)

func main() {

	wallet1 := wallet.NewWallet()
	wallet2 := wallet.NewWallet()

	transaction := wallet.NewTransaction(wallet1, wallet2.PublicKey, 60)
	transaction.Update(wallet1, wallet2.PublicKey, 100)

	fmt.Println(transaction)

	//httpPort := flag.String("http", "", "HTTP server port")
	//hostPort := flag.String("host", "", "WebSocket server port")
	//flag.Parse()
	//
	//p2p.Run(*httpPort, *hostPort)
	//
	//select {}
}
