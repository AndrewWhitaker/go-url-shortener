// @title        URL Shortener
// @version      1.0
// @description  A basic URL shortener

// @host      localhost:8080
// @BasePath  /api/v1
package main

import (
	"database/sql"
	"fmt"
	"url-shortener/db"
	"url-shortener/env"
	"url-shortener/jobs"
	"url-shortener/server"
	"url-shortener/services"
)

func main() {
	postgresHost := env.GetEnvVariable(env.PostgresHost)
	postgresPort := env.GetEnvVariable(env.PostgresPort)
	postgresUser := env.GetEnvVariable(env.PostgresUser)
	postgresPass := env.GetEnvVariable(env.PostgresPassword)
	postgresDatabase := env.GetEnvVariable(env.PostgresDatabase)

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
