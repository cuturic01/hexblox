package p2p

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hexblox/internal/domain"
	"net/http"
)

func SetBlockchainRoutes(node *Node) {
	blockchainGroup := node.HttpServer.Group("/hexblox")
	{
		blockchainGroup.GET("/blocks", func(context *gin.Context) {
			context.JSON(http.StatusOK, node.Blockchain.Chain())
		})

		blockchainGroup.POST("/mine", func(context *gin.Context) {
			var requestData struct {
				Data []*domain.Transaction `json:"data"`
			}

			if err := context.BindJSON(&requestData); err != nil {
				fmt.Println(err.Error())
				return
			}

			block := node.Blockchain.AddBlock(requestData.Data)
			fmt.Println(block)
			node.PropagateChain()
		})

		blockchainGroup.GET("/mine-transaction", func(context *gin.Context) {
			node.Mine()
			context.JSON(http.StatusOK, node.Blockchain.Chain())
		})
	}

}

func SetTransactionRoutes(node *Node) {
	transactionGroup := node.HttpServer.Group("hexblox/transactions")
	{
		transactionGroup.GET("", func(context *gin.Context) {
			context.JSON(http.StatusOK, node.TransactionPool.Transactions)
		})

		transactionGroup.POST("/transact", func(context *gin.Context) {
			var requestData struct {
				Recipient string  `json:"recipient"`
				Amount    float64 `json:"amount"`
			}
			if err := context.BindJSON(&requestData); err != nil {
				fmt.Println(err.Error())
				return
			}
			transaction := node.Wallet.CreateTransaction(
				requestData.Recipient,
				requestData.Amount,
				node.TransactionPool,
				node.Blockchain,
			)
			node.PropagateTransaction(transaction)
		})

		transactionGroup.GET("/public-key", func(context *gin.Context) {
			context.JSON(http.StatusOK, node.Wallet.PublicKey)
		})
	}
}
