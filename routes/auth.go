package routes

import (
	controller "github.com/RAVAN0407/jwt-token-auth/controllers"

	"github.com/gin-gonic/gin"
)

func Auth(routes *gin.Engine) {
	routes.POST("/user/signup", controller.SignUp())
	routes.POST("/user/login", controller.Login())
}
