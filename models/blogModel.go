package models

import (
	"akh/blog/config"
	"akh/blog/requests"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Blog struct {
	// Id          primitive.ObjectID `json:"_id"`
	Title       string             `json:"title"`
	Image       string             `json:"image"`
	Description string             `json:"description"`
	Tag         string             `json:"tags"`
	Comments    []bson.M           `json:"comments"`
	CreatedAt   primitive.DateTime `json:"created_at"`
	UpdatedAt   primitive.DateTime `json:"updated_at"`
}

var blogCollection *mongo.Collection = config.GetCollection(config.DB, "blog")

func AddBlog(requests requests.BlogRequest) (interface{}, bool) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	blogs := Blog{
		Title:       requests.Title,
		Image:       "",
		Description: requests.Description,
		Tag:         requests.Tags,
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now().UTC()), // DateTime generated from current Time
	}

	result, err := blogCollection.InsertOne(ctx, blogs)
	if err != nil {
		log.Fatal(err.Error())
		return err, false
	}

	return result, true
}
