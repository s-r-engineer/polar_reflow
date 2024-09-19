package mygin

import "github.com/gin-gonic/gin"

func Run() {
	engine := gin.Default()
	engine.PUT("/upload_data", func(ctx *gin.Context) {})
	engine.GET("/hrv/rmssi", func(ctx *gin.Context) {})
	engine.Run("localhost:6969")
}
