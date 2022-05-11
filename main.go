package main

import (
	"database/sql"
	"fmt"
	"os"
	"url-shortener/server"
)

func main() {
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPass := os.Getenv("POSTGRES_PASSWORD")
	postgresDatabase := os.Getenv("POSTGRES_DATABASE")

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		postgresUser,
		postgresPass,
		postgresHost,
		postgresPort,
		postgresDatabase,
	)

	fmt.Println(url)

	fmt.Printf("Connecting to %s\n", url)
	db, err := sql.Open("pgx", url)

	if err != nil {
		panic(fmt.Sprintf("Unable to connect to postgres: %s", err))
	}

	config := server.ServerConfig{DB: db}
	server.SetupServer(&config).Run()
}
