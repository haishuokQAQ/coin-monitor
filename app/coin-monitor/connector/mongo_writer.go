package connector

import (
	"context"

	"github.com/haishuokQAQ/coin-monitor/app/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func InitMongoConnection() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoDBUri))
	if err != nil {
		panic(err)
	}
	mongoClient = client
}

func GetMongoConnection() *mongo.Client {
	return mongoClient
}
