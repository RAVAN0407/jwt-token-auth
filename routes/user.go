package routes

import (
	controller "jwt-token-auth/controllers"
	"jwt-token-auth/middleware"

	"github.com/gin-gonic/gin"
)

func User(routes *gin.Engine) {
	routes.Use(middleware.Authenticate())
	routes.GET("/users", controller.GetUsers())
	routes.GET("/users/:user_id", controller.GetUsers())
}
