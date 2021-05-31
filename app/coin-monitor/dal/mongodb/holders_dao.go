package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/connector"
	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/model/mongo"
	"github.com/haishuokQAQ/coin-monitor/app/config"
)

func InsertHolderInfo(ctx context.Context, timeRange string, info *mongo.ERC20TokenHolderInf) error {
	collection := connector.GetMongoConnection().
		Database(config.CoinHolderDatabaseName).
		Collection(fmt.Sprintf("%s_%s", info.BaseCollectionName(), timeRange))
	_, err := collection.InsertOne(ctx, info)
	if err != nil {
		return err
	}
	return nil
}

func GetCurrentHolderInfoForAddress(ctx context.Context, address string) (*mongo.ERC20TokenHolderInf, error) {
	collection := connector.GetMongoConnection().
		Database(config.CoinHolderDatabaseName).
		Collection(fmt.Sprintf("%s_%s", mongo.ERC20TokenHolderInf{}.BaseCollectionName(), "current"))
	result := collection.FindOne(ctx, bson.D{primitive.E{
		Key:   "address",
		Value: address,
	}})
	if result.Err() != nil {
		return nil, result.Err()
	}
	resultInf := &mongo.ERC20TokenHolderInf{}
	err := result.Decode(resultInf)
	if err != nil {
		return nil, err
	}
	return resultInf, nil
}
