package back

import (
	"akh/blog/config"
	"akh/blog/helpers"
	"akh/blog/models"
	"akh/blog/requests"
	"akh/blog/responses"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var validate = validator.New()
var userCollection *mongo.Collection = config.GetCollection(config.DB, "user")

// var sampleSecretKey = []byte("SecretYouShouldHide")

// type Claims struct {
// 	Username string `json:"username"`
// 	jwt.RegisteredClaims
// }

// User Login
func Login(c *gin.Context) {

	var opt options.FindOneOptions
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data requests.SignInRequest

	// var tokenNUser responses.SignDataResponse

	var userData models.User

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, responses.SignInResponse{Status: http.StatusBadRequest, Message: err.Error(), Data: map[string]interface{}{"data": data}})
	}

	if validationErr := validate.Struct(&data); validationErr != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
		return
	}

	// opt.SetProjection(bson.M{"id": 1, "name": 1})
	err := userCollection.FindOne(ctx, bson.M{"email": data.Username}, &opt).Decode(&userData)

	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: err.Error(), Data: map[string]interface{}{}})
		return
	}

	if !helpers.CheckPasswordHash(data.Password, userData.Password) {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "password", Data: map[string]interface{}{"data": "Password does not match."}})
		return
	}

	// User Token GenerateToken
	// Returns true if the request was successful.
	tokenString, err := models.GenerateToken(userData)

	if err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: err.Error(), Data: map[string]interface{}{}})
		return
	}

	c.JSON(http.StatusOK, responses.SignInResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"token": tokenString, "id": userData.Id, "name": userData.Name, "username": userData.Email}})
}

// User Register
func Cteate(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var data requests.SignUpRequest

	defer cancel()

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if validationErr := validate.Struct(&data); validationErr != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
		return
	}

	hash, _ := helpers.HashPassword(data.Password)
	newUser := models.User{
		Id:        primitive.NewObjectID(),
		Name:      data.Name,
		Password:  hash,
		Email:     data.Email,
		CreatedAt: time.Now().In(*&time.Local),
	}

	results, err := userCollection.InsertOne(ctx, newUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: err.Error(), Data: map[string]interface{}{"data": err.Error()}})
	}

	c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": results}})
}
