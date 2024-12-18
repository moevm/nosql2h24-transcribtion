package db

import (
	"context"
	"fmt"
	"github.com/moevm/nosql2h24-transcribtion/config"
	"github.com/moevm/nosql2h24-transcribtion/models"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"os"
)

import (
	"go.mongodb.org/mongo-driver/bson"
)

func loadDataFromFile(filePath string, result interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = bson.UnmarshalExtJSON(data, true, result)
	if err != nil {
		return err
	}
	return nil
}

func SeedData(cfg config.Config, client *mongo.Client) {

	ctx := context.Background()

	var users []models.User
	var jobs []models.Job
	var servers []models.Server

	if err := loadDataFromFile("db/seed_data/users.json", &users); err != nil {
		log.Fatal("Error loading users data: ", err)
	}

	if err := loadDataFromFile("db/seed_data/jobs.json", &jobs); err != nil {
		log.Fatal("Error loading jobs data: ", err)
	}

	if err := loadDataFromFile("db/seed_data/servers.json", &servers); err != nil {
		log.Fatal("Error loading servers data: ", err)
	}

	usersCollection := client.Database(cfg.DBName).Collection("users")
	jobsCollection := client.Database(cfg.DBName).Collection("jobs")
	serversCollection := client.Database(cfg.DBName).Collection("servers")

	_, err := usersCollection.DeleteMany(ctx, map[string]interface{}{})
	if err != nil {
		log.Fatal("Error deleting users data: ", err)
	}
	_, err = jobsCollection.DeleteMany(ctx, map[string]interface{}{})
	if err != nil {
		log.Fatal("Error deleting jobs data: ", err)
	}
	_, err = serversCollection.DeleteMany(ctx, map[string]interface{}{})
	if err != nil {
		log.Fatal("Error deleting servers data: ", err)
	}

	for _, server := range servers {
		_, err := serversCollection.InsertOne(ctx, server)
		if err != nil {
			log.Fatal("Error inserting server: ", err)
		}
	}

	for _, job := range jobs {
		_, err := jobsCollection.InsertOne(ctx, job)
		if err != nil {
			log.Fatal("Error inserting job: ", err)
		}
	}

	for _, user := range users {
		_, err := usersCollection.InsertOne(ctx, user)
		if err != nil {
			log.Fatal("Error inserting user: ", err)
		}
	}

	fmt.Println("Data seeded successfully!")
}
