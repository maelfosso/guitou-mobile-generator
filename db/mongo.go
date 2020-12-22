package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// // MongoDBClient store the connexion to MongoDB instance
// type MongoDBClient struct {
// 	client  *mongo.Client
// 	dbname  string
// 	context context.Context
// }

// var _mongo *MongoDBClient

// MongoDBClient store the connexion to MongoDB instance
type MongoDBClient struct {
	uri    string
	dbname string

	client  *mongo.Client
	context context.Context
}

// NewMongoDBClient create a struct
func NewMongoDBClient(uri, dbname string) *MongoDBClient {
	return &MongoDBClient{
		uri,
		dbname,

		nil,
		nil,
	}
}

// Connect to connect to the MongoDB database dbname
func (mdbc *MongoDBClient) Connect() {
	clientOptions := options.Client().ApplyURI(mdbc.uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	ctxPing, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctxPing, readpref.Primary())
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	log.Println("Connected to MongoDB")
	// _mongo = &MongoDBClient{
	// 	client,
	// 	dbname,
	// 	ctx,
	// }

	mdbc.client = client
	mdbc.context = ctx
}

// GetCollection eases the access to the Database collection
func (mdbc *MongoDBClient) GetCollection(collection string) *mongo.Collection {
	// return _mongo.client.Database(_mongo.dbname).Collection(collection)
	return mdbc.client.Database(mdbc.dbname).Collection(collection)
}

// Close closes the connection to MongoDB instance
func (mdbc *MongoDBClient) Close() {
	// err := _mongo.client.Disconnect(_mongo.context)
	err := mdbc.client.Disconnect(mdbc.context)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	log.Println("Connection to MongoDB closed")
}
