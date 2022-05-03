package main

import (
	"context"
	"testing"
)

func TestMain(t *testing.T) {
	ctx := context.Background()

	container, db, err := CreateTestContainer(ctx, "testdb")

	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()
	defer container.Terminate(ctx)
}
