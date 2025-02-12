package integrationtest

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/braumsmilk/go-db"
	"github.com/braumsmilk/go-db/testutils"
)

func TestMain(m *testing.M) {
	// Initialize the Redis container
	redisContainer := testutils.InitRedis()
	defer func() {
		if err := redisContainer.Terminate(context.Background()); err != nil {
			log.Fatalf("failed to terminate redis container: %v", err)
		}
	}()

	// Initialize the Postgres container
	postgresContainer := testutils.InitPostgres(nil, nil)
	defer func() {
		if err := postgresContainer.Terminate(context.Background()); err != nil {
			log.Fatalf("failed to terminate postgres container: %v", err)
		}
	}()

	os.Exit(m.Run())
}

func TestConnect(t *testing.T) {
	// This is a placeholder test to ensure that the test suite runs correctly.
	// You can add your actual test logic here.
	t.Log("Running integration tests...")

	db.InitDefault()
	
}
