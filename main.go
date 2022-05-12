package main

import (
	"database/sql"
	"fmt"
	"os"
	"url-shortener/db"
	"url-shortener/jobs"
	"url-shortener/server"
	"url-shortener/services"
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
	sqlDB, err := sql.Open("pgx", url)

	if err != nil {
		panic(fmt.Sprintf("Unable to connect to postgres: %s", err))
	}

	gormDB, err := db.ConnectDatabase(sqlDB)

	jobs.StartScheduler(gormDB, services.SystemClock{})

	config := server.ServerConfig{DB: gormDB}
	server.SetupServer(&config).Run()
}
