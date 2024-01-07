package connection

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Users    *mongo.Collection
	Posts    *mongo.Collection
	Comments *mongo.Collection
}

var db Database

func GetDB() Database {
	return db
}

func Init() *mongo.Client {
	MONGO_URI := os.Getenv("MONGO_URI")
	if MONGO_URI == "" {
		panic("MONGO_URI is not set")
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(MONGO_URI))

	if err != nil {
		panic(err)
	}

	ping(client)
	db = NewDatabase(client)
	return client
}

func ping(client *mongo.Client) {
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

func Disconnect(client *mongo.Client) {
	if err := client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}

func NewDatabase(client *mongo.Client) Database {
	return Database{
		Users:    client.Database("lorem_ipsum").Collection("users"),
		Posts:    client.Database("lorem_ipsum").Collection("posts"),
		Comments: client.Database("lorem_ipsum").Collection("comments"),
	}
}
