package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hexblox/internal/blockchain"
	"net/http"
)

func SetBlockchainRoutes(engine *gin.Engine, blockchain *blockchain.Blockchain[string]) {
	engine.GET("/blocks", func(context *gin.Context) {
		context.JSON(http.StatusOK, blockchain.Chain())
	})

	engine.POST("/mine", func(context *gin.Context) {
		var requestData struct {
			Data []*string `json:"data"`
		}

		if err := context.BindJSON(&requestData); err != nil {
			fmt.Println(err.Error())
			return
		}

		block := blockchain.AddBlock(requestData.Data)
		fmt.Println(block)
	})
}
