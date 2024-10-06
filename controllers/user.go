package controller

import (
	"context"
	"fmt"
	"jwt-token-auth/database"
	"jwt-token-auth/helpers"
	"jwt-token-auth/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.Client.OpenConnection("user")
var validate = validator.New()

func HashPassword(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(hash)
}

func VerifyPassword(stored, new string) error {
	err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(new))
	if err != nil {
		return fmt.Errorf("email or password is incorrect")
	}
	return nil
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err", err.Error()})
			return
		}
		err = validate.Struct(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err", err.Error()})
			return
		}
		ecount, err := userCollection.CountDocuments(ctx, bson.M{"email", user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		pcount, err := userCollection.CountDocuments(ctx, bson.M{"phone", user.Phone})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		if ecount > 0 || pcount > 0 {
			c.JSON(http.StatusInternalServerError, "email or phone already in use")
		}
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()

		insertValue, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, insertValue)

	}

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user, foundUser models.User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err", err.Error()})
			return
		}

		err = userCollection.FindOne(ctx, bson.M{"email", user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		err = VerifyPassword(*foundUser.Password, *user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "user not found"})
			return
		}

		tkn, refreshtkn, err := helper.GenerateAllToken(*foundUser.Email, *foundUser.FirstName,
			*foundUser.LastName, *foundUser.UserType, foundUser.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		helper.UpdateAllTokens(tkn, refreshtkn, foundUser.UserID)

		c.JSON(http.StatusOK, foundUser)

	}

}

func GetUsers() {
	return func(c *gin.Context) {
		helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error", err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		recordPerPage, err = strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}}, 
			{"total_count", bson.D{{"$sum", 1}}}, 
			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},}}}
		result,err := userCollection.Aggregate(ctx, mongo.Pipeline{
						matchStage, groupStage, projectStage})
		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing user items"})
		}
		var allusers []bson.M
		if err = result.All(ctx, &allusers); err!=nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allusers[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		if err := helpers.MatchUserTypeToUid(c, userID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id", userID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
