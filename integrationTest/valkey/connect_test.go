package valkey

import (
	"context"
	"log"
	"testing"

	"github.com/ooqls/go-db/containers"
	"github.com/ooqls/go-db/valkey"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Initialize the Valkey container
	valkeyContainer := containers.StartValkey(context.Background())
	defer func() {
		if err := valkeyContainer.Terminate(context.Background()); err != nil {
			log.Fatalf("failed to terminate valkey container: %v", err)
		}
	}()

	m.Run()
}

func TestConnectValkey(t *testing.T) {
	ctx := context.Background()
	con := valkey.GetConnection(ctx)
	assert.NotNilf(t, con, "valkey client should not be nil")
}
