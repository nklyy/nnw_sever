package mongodb

import (
	"context"
	"nnw_s/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewConn(cfg *config.Config) (*mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoDbUrl))
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		defer client.Disconnect(ctx)
		return nil, err
	}

	collection := client.Database(cfg.MongoDbName)

	return collection, nil
}
