package db

import (
	"context"
	"fmt"
	"github.com/moevm/nosql2h24-transcribtion/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const ConnectionTimeout = 10 * time.Second

func InitConnection(cfg *config.Config) *mongo.Client {

	ctx, cancel := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer cancel()

	mongoconn := options.Client().ApplyURI(cfg.DBUri)
	client, err := mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal(err)
	}

	//TODO Возможно реализовать отключение от бд

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}
