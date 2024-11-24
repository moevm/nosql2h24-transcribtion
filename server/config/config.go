package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DBUri  string `mapstructure:"MONGODB_LOCAL_URI"`
	Port   string `mapstructure:"PORT"`
	DBName string `mapstructure:"MONGODB_LOCAL_NAME"`
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return Config{
			DBUri:  os.Getenv("MONGODB_URI"),
			Port:   os.Getenv("PORT"),
			DBName: os.Getenv("MONGODB_NAME"),
		},
		nil
}
