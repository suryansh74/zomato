package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func DBConnect(mongoURI string) *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	// context only for the initial connect + ping
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // ✅ always cancel to free resources, safe here since we only use it for ping

	client, err := mongo.Connect(opts)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to mongodb: %v", err))
	}

	// ping to confirm connection is alive
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(fmt.Sprintf("failed to ping mongodb: %v", err))
	}

	fmt.Println("✅ connected to mongodb atlas")
	return client // ✅ return client WITHOUT disconnecting
}
