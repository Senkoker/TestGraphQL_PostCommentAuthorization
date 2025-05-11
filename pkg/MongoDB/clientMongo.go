package MongoDB

import (
	"context"
	"friend_graphql/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type ClientMongo struct {
	Client *mongo.Collection
}

func NewClientMongo(cfg *config.Cfg) *ClientMongo {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.URL))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database(cfg.Mongo.Database).Collection(cfg.Mongo.CollectionName)
	return &ClientMongo{Client: collection}
}
