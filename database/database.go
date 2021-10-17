package database

import (
	"context"
	"fmt"
	"log"
	"skillshare/converter/config"

	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func Init() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	databaseConfig := config.GetDatabaseConfig()
	url := fmt.Sprintf("%s", databaseConfig.URL)
	tmpClient, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}
	client = tmpClient
	log.Println("DB Connected")
}

func GetDatabaseClient() *mongo.Client {
	if client == nil {
		Init()
	}
	return client
}

func Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
	log.Println("DB Disconnected")
}
