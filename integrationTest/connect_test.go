package integrationtest

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/ooqls/go-db/pgx"
	"github.com/ooqls/go-db/sqlx"
	"github.com/ooqls/go-db/testutils"
	"github.com/stretchr/testify/assert"
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
	postgresContainer := testutils.StartPostgres(context.Background())
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

	assert.Nilf(t, sqlx.InitDefault(), "Expected InitDefault to return nil, but got an error.")

	assert.Nilf(t, pgx.InitDefault(), "Expected InitDefault to return nil, but got an error.")

}
