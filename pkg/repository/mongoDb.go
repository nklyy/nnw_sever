package repository

import (
	"context"
	"log"
	"nnw_s/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func MongoDbConnection(cfg *config.Configurations) (*mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoDbUrl))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		defer client.Disconnect(ctx)
		log.Fatal(err)
	}

	collection := client.Database("NNW")

	return collection, nil
}
