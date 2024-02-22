package main

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongodbURI := os.Getenv("MONGODB_URI")
	if mongodbURI == "" {
		panic("Specify mongodb connection URI")
	}
	mongodbDBName := os.Getenv("MONGODB_DB")
	if mongodbDBName == "" {
		panic("Specify mongodb database connection!")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	serverOptions := options.Client().ApplyURI(mongodbURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, serverOptions)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	var result bson.M
	if err = client.Database(mongodbDBName).RunCommand(ctx, bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Succesful connection!")
}
