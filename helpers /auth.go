package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) error {
	userType := c.GetString("user_type")
	if userType != role {
		return errors.New("Unauthorized Access")
	}
	return nil
}

func MatchUserTypeToUid(c *gin.Context, userID string) error {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")

	if userType == "USER" && uid != userID {
		return errors.New("Unauthorised Access")
	}
	return CheckUserType(c, userType)
}
