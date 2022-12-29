package models

import (
	"akh/blog/config"
	"akh/blog/helpers"
	"akh/blog/requests"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Id        primitive.ObjectID `json:"id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	Password  string             `json:"password"`
	CreatedAt time.Time          `json:"created_at"`
}

var userCollection *mongo.Collection = config.GetCollection(config.DB, "user")

func AuthRequest() func(userReq requests.SignInRequest) interface{} {
	return func(userReq requests.SignInRequest) interface{} {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user User

		err := userCollection.FindOne(ctx, bson.M{"name": userReq.Username}).Decode(&user)

		if err != nil {
			log.Fatal(err)
			return "errors"
		}

		if helpers.CheckPasswordHash(userReq.Password, user.Password) {
			fmt.Println(user)
		}

		return user
	}
}

type Claims struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Id       string `json:"id"`
	jwt.RegisteredClaims
}

var sampleSecretKey = []byte("SecretYouShouldHide")

func GenerateToken(userData User) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Id:       userData.Id.String(),
		Username: userData.Email,
		Name:     userData.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(sampleSecretKey)
	return tokenString, err
}
