package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

var RecipesCollection *mongo.Collection

func ConnectDB() {
	str := viper.GetString("db_uri")
	fmt.Println(str)
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("db_uri")))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // TODO: Length?
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	Client = client
	RecipesCollection = Client.Database("soups-up").Collection("recipes")
}
