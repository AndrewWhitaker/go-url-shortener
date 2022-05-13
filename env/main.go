package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	PostgresHost     = "POSTGRES_HOST"
	PostgresPort     = "POSTGRES_PORT"
	PostgresUser     = "POSTGRES_USER"
	PostgresPassword = "POSTGRES_PASSWORD"
	PostgresDatabase = "POSTGRES_DATABASE"

	GinMode = "GIN_MODE"
)

func GetEnvVariable(key string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("error getting environment variable %s: %v", key, err)
	}

	return os.Getenv(key)
}
