package elasticsearch

import (
	"context"
	"testing"

	"github.com/ooqls/go-db/testutils"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	ctx := context.Background()

	container := testutils.StartElasticsearch(ctx)
	defer container.Terminate(ctx)

	// Get the mapped port
	err := InitDefault()
	assert.Nilf(t, err, "failed to initialize elasticsearch client: %v", err)

	l.Info("Elasticsearch client initialized successfully")
}
