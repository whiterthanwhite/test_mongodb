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
	dbName := os.Getenv("MONGODB_DB")
	if dbName == "" {
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

	if err = testDatabaseConnection(ctx, client, dbName); err != nil {
		panic(err)
	}

	if err = testInsertOne(ctx, client, dbName, "persons"); err != nil {
		panic(err)
	}

	if err = testInsertMany(ctx, client, dbName, "persons"); err != nil {
		panic(err)
	}

	if err = testFindOne(ctx, client, dbName, "persons"); err != nil {
		panic(err)
	}

	if err = testFind(ctx, client, dbName, "persons"); err != nil {
		panic(err)
	}

	if err = testDeleteMany(ctx, client, dbName, "persons"); err != nil {
		panic(err)
	}
}

func testDatabaseConnection(ctx context.Context, client *mongo.Client, dbName string) error {
	var result bson.M
	if err := client.Database(dbName).RunCommand(ctx, bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return err
	}
	fmt.Println("Succesful connection!")
	return nil
}

type Person struct {
	FirstName string `bson:"first_name,omitempty"`
	LastName  string `bson:"last_name,omitempty"`
}

func testInsertOne(ctx context.Context, client *mongo.Client, dbName, collName string) error {
	// if collection doesn't exist. it will be created automatically
	coll := client.Database(dbName).Collection(collName)
	p := &Person{
		FirstName: "Test Person First Name",
		LastName:  "Test Person Last Name",
	}

	_, err := coll.InsertOne(ctx, p)
	if err != nil {
		return err
	}

	fmt.Println("Person", *p, "created successfully")
	return nil
}

func testInsertMany(ctx context.Context, client *mongo.Client, dbName, collName string) error {
	coll := client.Database(dbName).Collection(collName)
	persons := []interface{}{
		&Person{
			FirstName: "Leonardo",
			LastName:  "Turtle",
		},
		&Person{
			FirstName: "Donatello",
			LastName:  "Turtle",
		},
	}

	_, err := coll.InsertMany(ctx, persons)
	if err != nil {
		return err
	}

	fmt.Println(len(persons), "documents created successfully")
	return nil
}

func testFindOne(ctx context.Context, client *mongo.Client, dbName, collName string) error {
	coll := client.Database(dbName).Collection(collName)
	personFilter := bson.D{{"first_name", "Donatello"}}

	var p Person
	err := coll.FindOne(ctx, personFilter).Decode(&p)
	if err != nil {
		return err
	}

	fmt.Println("Find", p)
	return nil
}

func testFind(ctx context.Context, client *mongo.Client, dbName, collName string) error {
	coll := client.Database(dbName).Collection(collName)
	personFilter := bson.D{{"first_name", "Donatello"}}

	cur, err := coll.Find(ctx, personFilter)
	if err != nil {
		return err
	}

	var persons []Person
	if err = cur.All(ctx, &persons); err != nil {
		return err
	}

	fmt.Println("Retreived", len(persons))
	return nil
}

func testDeleteMany(ctx context.Context, client *mongo.Client, dbName, collName string) error {
	coll := client.Database(dbName).Collection(collName)
	filter := bson.D{{}}

	res, err := coll.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	fmt.Println("Deleted", res.DeletedCount, "documents!")
	return nil
}
