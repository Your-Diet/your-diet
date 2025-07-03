package utils

import (
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	MongoURL string
	DBName   string
	Port     string
}

func LoadEnvConfig() *EnvConfig {
	scope := os.Getenv("SCOPE")

	if scope == "" || scope != "prod" {
		err := godotenv.Load(".env")

		if err != nil {
			panic("Error loading .env file: " + err.Error())
		}
	}

	mongoURL := os.Getenv("MONGODB_URL")
	if mongoURL == "" {
		panic("MONGODB_URL is not set")
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		panic("MONGO_DB_NAME is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is not set")
	}

	return &EnvConfig{
		MongoURL: mongoURL,
		DBName:   dbName,
		Port:     port,
	}
}
