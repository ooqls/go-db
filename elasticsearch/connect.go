package elasticsearch

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/ooqls/go-db/dbconfig"
	"github.com/ooqls/go-log"
	"github.com/ooqls/go-registry"
	"go.uber.org/zap"
)

var l *zap.Logger = log.NewLogger("elasticsearch")
var esClient *elasticsearch.TypedClient
var m sync.Mutex

func Init(db *registry.Database) error {
	m.Lock()
	defer m.Unlock()

	cfg, err := dbconfig.Resolve(db)
	if err != nil {
		return fmt.Errorf("failed to resolve elasticsearch config: %v", err)
	}

	// headers := http.Header{}
	// headers.Del("Accept")
	// headers.Set("Accept", "application/json")
	// headers.Set("Content-Type", "application/json")
	esClient, err = elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("https://%s", cfg.Addr()),
		},
		Username:          cfg.Username,
		Password:          cfg.Password,
		EnableDebugLogger: true,
		Transport: &http.Transport{
			TLSClientConfig: cfg.TLS,
		},
		Logger: &elastictransport.TextLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		},
	})
	if err != nil {
		return fmt.Errorf("error creating the client: %s", err)
	}
	ctx := context.Background()
	success := false
	for i := 0; i < 5; i++ {
		if i > 0 {
			l.Info("Retrying to connect to elasticsearch", zap.Int("attempt", i))
			time.Sleep(time.Second)
		}

		info, err := esClient.Info().Do(ctx)
		if err != nil {
			l.Error(fmt.Sprintf("failed to get elasticsearch info, attempt %d", i), zap.Error(err))
			continue
		}

		l.Debug("Connected to Elasticsearch", zap.String("version", info.Version.Int))

		l.Debug("Elasticsearch client initialized successfully")
		success = true
		break
	}

	if !success {
		return fmt.Errorf("failed to initialize elasticsearch client")
	}

	return nil
}

func InitDefault() error {
	reg := registry.Get()
	if reg.Elasticsearch == nil {
		return fmt.Errorf("elasticsearch not found in registry")
	}

	return Init(reg.Elasticsearch)
}

func Get() *elasticsearch.TypedClient {
	m.Lock()
	defer m.Unlock()

	return esClient
}
