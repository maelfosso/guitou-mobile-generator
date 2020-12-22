package db

import (
	"os"
)

// Init initializes the database
func Init() *MongoDBClient {
	host := os.Getenv("MONGO_HOST")
	dbName := os.Getenv("MONGO_DBNAME")

	mdbc := NewMongoDBClient(host, dbName)
	mdbc.Connect()

	return mdbc
}
