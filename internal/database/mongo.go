package database 

import (
	"context"
	"time"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/solonode/golang-jwt-mongo/config"
)

func SetupMongoDB() *mongo.Client{	
	// Loading env vars
	config, err := config.LoadENVFile("../../")  // file path of app.env
	if err != nil{
		log.Fatal("Error loading env vars: ", err)
	}

	// context with timeout
	ctx, cancel :=  context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.LocalDBUri)  //mongodb connection options

	// mongodb connection
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil{
		return nil
	}
	
	// ping the mongodb server to check connection
	if err := client.Ping(ctx, nil); err != nil{
		return nil
	}
	
	log.Printf("MongoDB connected successfully on: %v", config.LocalDBUri)
	return client
}

var DB *mongo.Client = SetupMongoDB()   // mongodb client instance


// Get all collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {

	// Loading env vars
	config, err := config.LoadENVFile("../../")  // file path of app.env
	if err != nil{
		log.Fatal("Error loading env vars: ", err)
	}

	collection := client.Database(config.DBName).Collection(collectionName)
	return collection
}

