package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	User      *mongo.Collection
	Frequnecy *mongo.Collection
	Payloads  *mongo.Collection
}

// Connect to mongodb server and returns type {*Mongodb} if successfully connected.
func ConnectMongoDB(env *ENV) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var mongoClient *mongo.Client
	var err error

	if env.MONGODB_CONNECTION_METHOD == "manual" {
		credential := options.Credential{
			Username: env.MONGODB_USERNAME,
			Password: env.MONGODB_PASSWORD,
		}

		clientOptions := options.Client().ApplyURI("mongodb://" + env.MONGODB_HOST + ":" + env.MONGODB_PORT).SetAuth(credential)
		mongoClient, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			defer cancel()
			return nil, err
		}
	} else if env.MONGODB_CONNECTION_METHOD == "auto" {
		serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
		clientOptions := options.Client().ApplyURI(env.MONGODB_CONNECTION_STRING).SetServerAPIOptions(serverAPIOptions)
		mongoClient, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			defer cancel()
			return nil, err
		}
	}
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		defer cancel()
		return nil, err
	}

	var mongodb MongoDB

	storage := mongoClient.Database("storage")
	mongodb.User = storage.Collection("user")
	mongodb.Frequnecy = storage.Collection("userFrequencyTable")

	delivery := mongoClient.Database("delivery")
	mongodb.Payloads = delivery.Collection("payloads")

	defer cancel()
	return &mongodb, nil
}
