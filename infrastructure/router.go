package infrastructure

import (
	"app/interfaces/controllers"
	"log"

	"github.com/gin-gonic/gin"
)

func Router() {
	r := gin.Default()

	BotApiController := controllers.NewBotApiController(Init())

	// ヘルスチェック用エンドポイント
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	r.POST("/bots", func(c *gin.Context) { BotApiController.Create(c) })
	r.GET("/bots", func(c *gin.Context) { BotApiController.Index(c) })
	r.GET("/bot/:id", func(c *gin.Context) { BotApiController.Show(c) })

	// HTTPサーバーを非同期で起動
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatal("Error starting HTTP server: ", err)
		}
	}()
	log.Printf("Start")

}
