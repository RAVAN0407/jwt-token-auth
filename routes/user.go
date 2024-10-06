package routes

import (
	controller "github.com/RAVAN0407/jwt-token-auth/controllers"
	middleware "github.com/RAVAN0407/jwt-token-auth/middleware"

	"github.com/gin-gonic/gin"
)

func User(routes *gin.Engine) {
	routes.Use(middleware.Authenticate())
	routes.GET("/users", controller.GetUsers())
	routes.GET("/users/:user_id", controller.GetUsers())
}
