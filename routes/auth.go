package routes

import (
	controller "jwt-token-auth/controllers"

	"github.com/gin-gonic/gin"
)

func Auth(routes *gin.Engine) {
	routes.POST("/user/signup", controller.SignUp())
	routes.POST("/user/login", controller.Login())
}
