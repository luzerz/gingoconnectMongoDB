package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		fmt.Print("Error loading .env file then use env server")
	}

	return os.Getenv("MONGOURI")
}
