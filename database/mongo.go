package database

import (
	"context"
	"cowingram/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MongoBotClient struct {
	Client *mongo.Client       `json:"client,omitempty"`
	Config *config.MongoConfig `json:"config,omitempty"`
}

func (Bc *MongoBotClient) GetClient() {

	Bc.Config.Init()

	mongoUrl, err := Bc.Config.BuildUrl()

	if err != nil {
		log.Panic(err)
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUrl))

	if err != nil {
		log.Panic(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Panic(err)
		}
	}(client, ctx)

	log.Println("mongo-db connection completed.")
}

func Connect() {
	var mgc MongoBotClient
	mgc.GetClient()
}
