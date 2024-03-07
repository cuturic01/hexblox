package p2p

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetBlockchainRoutes(node *Node) {
	node.HttpServer.GET("/blocks", func(context *gin.Context) {
		context.JSON(http.StatusOK, node.Blockchain.Chain())
	})

	node.HttpServer.POST("/mine", func(context *gin.Context) {
		var requestData struct {
			Data []*string `json:"data"`
		}

		if err := context.BindJSON(&requestData); err != nil {
			fmt.Println(err.Error())
			return
		}

		block := node.Blockchain.AddBlock(requestData.Data)
		fmt.Println(block)
		node.PropagateChain()
	})
}
