package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	database "github.com/RAVAN0407/jwt-token-auth/database"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	firstName string
	lastName  string
	email     string
	userType  string
	uid       string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

var userCollection *mongo.Collection = database.Client.OpenConnection("user")

func GenerateAllToken(email string, firstName string, lastName string, userType string, uid string) (string, string, error) {
	claims := &SignedDetails{
		firstName: firstName,
		lastName:  lastName,
		email:     email,
		userType:  userType,
		uid:       uid,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return token, refreshToken, err
	}

	return token, refreshToken, err

}

func UpdateAllTokens(signedTkn, signedRefreshTkn string, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", signedTkn})
	updateObj = append(updateObj, bson.E{"refreshToken", signedRefreshTkn})
	updateObj = append(updateObj, bson.E{"updatedAt", time.Now()})

	upsert := true
	filter := bson.M{"userID": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	err := userCollection.UpdateOne{
		ctx,
		filter,
		bson.D{{"$set", updateObj}},
		&opt,
	}
	if err != nil {
		return err
	}
	return nil
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}
	return claims, msg
}
