package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MySQLAddr     string
	MySQLUser     string
	MySQLPassword string
	MySQLDatabase string

	RedisAddr     string
	RedisPassword string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return Config{
		MySQLAddr:     os.Getenv("MYSQL_ADDR"),
		MySQLUser:     os.Getenv("MYSQL_USER"),
		MySQLPassword: os.Getenv("MYSQL_PASSWORD"),
		MySQLDatabase: os.Getenv("MYSQL_DATABASE"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}
}
