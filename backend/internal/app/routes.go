package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerRoutes(router *gin.Engine, dependencies dependencies) {
	router.GET("/health", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
