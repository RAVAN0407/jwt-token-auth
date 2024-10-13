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
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct {
	FirstName string
	LastName  string
	Email     string
	UserType  string
	UserID    string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

var userCollection *mongo.Collection = database.OpenConnection(database.Client, "user")

func GenerateAllToken(email string, firstName string, lastName string, userType string, uid string) (string, string, error) {
	claims := &SignedDetails{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		UserType:  userType,
		UserID:    uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
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

	updateObj := bson.D{
		{"token", signedTkn},
		{"refreshToken", signedRefreshTkn},
		{"updatedAt", time.Now()},
	}

	filter := bson.M{"userID": userId}
	update := bson.M{"$set": updateObj}

	_, err := userCollection.UpdateOne(ctx, filter, update)
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
