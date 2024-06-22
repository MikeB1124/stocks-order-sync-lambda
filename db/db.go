package db

import (
	"context"
	"fmt"
	"time"

	"github.com/MikeB1124/stocks-order-sync-lambda/configuration"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func init() {
	config := configuration.GetConfig()
	// Connect to MongoDB
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.du0vf.mongodb.net", config.MongoDB.Username, config.MongoDB.Password))
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	mongoClient = client
}

func UpateOrder(order alpaca.Order) (*mongo.UpdateResult, error) {
	collection := mongoClient.Database("Stocks").Collection("orders")
	filter := bson.M{"order.id": order.ID}
	update := bson.M{
		"$set": bson.M{
			"order":           order,
			"recordUpdatedAt": time.Now().UTC(),
		},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return updateResult, nil
}
