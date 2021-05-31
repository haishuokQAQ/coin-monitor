package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/connector"
	"github.com/haishuokQAQ/coin-monitor/app/coin-monitor/model/mongo"
)

func WriteHolderAnalyserConfig(ctx context.Context, config *mongo.HolderAnalyseConfig) error {
	config, err := ReadHolderAnalyseConfig(ctx)
	if err != nil {
		return err
	}
	if config != nil {
		_, err = connector.GetMongoConnection().Database("coin_holder_data").Collection("config").UpdateOne(ctx, bson.D{primitive.E{
			Key:   "key",
			Value: "holder_analyse_config",
		}}, config)
	} else {
		_, err = connector.GetMongoConnection().Database("coin_holder_data").Collection("config").InsertOne(ctx, config)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadHolderAnalyseConfig(ctx context.Context) (*mongo.HolderAnalyseConfig, error) {
	result := connector.GetMongoConnection().Database("coin_holder_data").Collection("config").FindOne(ctx, bson.D{primitive.E{
		Key:   "key",
		Value: "holder_analyse_config",
	}})
	if err := result.Err(); err != nil {
		return nil, err
	}
	resultConfig := &mongo.HolderAnalyseConfig{}
	raw, err := result.DecodeBytes()
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, nil
	}
	err = result.Decode(resultConfig)
	if err != nil {
		return nil, err
	}
	return resultConfig, nil
}
