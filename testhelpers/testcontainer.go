/* https://brietsparks.com/testcontainers-golang-db-access/ */
package testhelpers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func CreateTestContainer(ctx context.Context, dbname string) (testcontainers.Container, *sql.DB, error) {
	var env = map[string]string{
		"POSTGRES_PASSWORD": "postgres",
		"POSTGRES_USER":     "postgres",
		"POSTGRES_DB":       dbname,
	}

	dbUrl := func(port nat.Port) string {
		return fmt.Sprintf("postgres://postgres:postgres@localhost:%s/%s?sslmode=disable", port.Port(), dbname)
	}

	port := "5432/tcp"

	req := testcontainers.ContainerRequest{
		Image:        "postgres:14.2",
		ExposedPorts: []string{port},
		Cmd:          []string{"postgres", "-c", "fsync=off"},
		Env:          env,
		WaitingFor:   wait.ForSQL(nat.Port(port), "postgres", dbUrl).Timeout(time.Second * 7),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return container, nil, fmt.Errorf("failed to start container: %s", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(port))

	if err != nil {
		return container, nil, fmt.Errorf("failed to get container external port: %s", err)
	}

	url := fmt.Sprintf("postgres://postgres:postgres@localhost:%s/%s?sslmode=disable", mappedPort.Port(), dbname)

	db, err := sql.Open("pgx", url)

	if err != nil {
		return container, db, fmt.Errorf("failed to establish database connection: %s", err)
	}

	return container, db, nil
}
