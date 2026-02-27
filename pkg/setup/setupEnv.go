package setup

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func ApiToken() string {
	return ReadEnv("API_TOKEN")
}

func ReadEnv(name string) string {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env: %v", err)
		return ""
	}

	env := os.Getenv(name)
	log.Printf("%s env: %s", name, env)

	return env
}
