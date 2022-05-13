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
	// https://github.com/joho/godotenv/issues/99
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalf("error getting environment variable %s: %v\n", key, err)
	}

	return os.Getenv(key)
}
