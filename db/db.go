package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var RecipesCollection *mongo.Collection

var Test string

func ConnectDB() {
	db_uri, ok := os.LookupEnv("DB_URI")
	if !ok {
		db_uri = "mongodb://localhost:27017/soups-up"
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(db_uri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
	RecipesCollection = client.Database("soups-up").Collection("recipes")
}
