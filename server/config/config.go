package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	DBUri        string `mapstructure:"MONGODB_LOCAL_URI"`
	Port         string `mapstructure:"PORT"`
	DBName       string `mapstructure:"MONGODB_LOCAL_NAME"`
	SeedDatabase bool   `mapstructure:"SEED_DATABASE"`
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	seedDatabase, err := strconv.ParseBool(os.Getenv("SEED_DATABASE"))
	if err != nil {
		log.Fatal("Error parsing SEED_DATABASE")
	}
	return Config{
			DBUri:        os.Getenv("MONGODB_URI"),
			Port:         os.Getenv("PORT"),
			DBName:       os.Getenv("MONGODB_NAME"),
			SeedDatabase: seedDatabase,
		},
		nil
}
