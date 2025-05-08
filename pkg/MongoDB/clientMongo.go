package MongoDB

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type ClientMongo struct {
	Client *mongo.Collection
}

func NewClientMongo() *ClientMongo {
	uri := "mongodb://admin:12345@localhost:27017"
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("postComment").Collection("postComment")
	return &ClientMongo{Client: collection}
}
