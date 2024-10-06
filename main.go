package main

import (
	routes "github.com/RAVAN0407/jwt-token-auth/routes"

	helpers "github.com/RAVAN0407/jwt-token-auth/helpers"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.New()
	router.Use(gin.Logger)
	router.Use(gin.Recovery)
	routes.Auth(router)
	routes.User(router)

	router.GET("/api-1", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"success": "Access granted"})

	})
	router.POST("/api-2", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"success": "Access granted"})

	})
	router.Run(":" + helpers.GetPort())
}
