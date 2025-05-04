package elasticsearch

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/ooqls/go-log"
	"github.com/ooqls/go-registry"
	"go.uber.org/zap"
)

var l *zap.Logger = log.NewLogger("elasticsearch")
var esClient *elasticsearch.TypedClient
var m sync.Mutex

type ElasticsearchOptions struct {
	Host            string
	Port            int
	User            string
	Pw              string
	DB              string
	InsecureSkipTLS bool
}

func Init(opts ElasticsearchOptions) error {
	m.Lock()
	defer m.Unlock()

	var err error
	trans := http.DefaultTransport.(*http.Transport)
	trans.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: opts.InsecureSkipTLS,
	}

	esClient, err = elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("https://%s:%d", opts.Host, opts.Port)},
		Username:  opts.User,
		Password:  opts.Pw,
		Transport: trans,
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
		
		_, err := esClient.Info().Do(ctx)
		if err != nil {
			l.Error(fmt.Sprintf("failed to get elasticsearch info, attempt %d", i), zap.Error(err))
			continue
		}

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

	opts := ElasticsearchOptions{
		Host:            reg.Elasticsearch.Host,
		Port:            reg.Elasticsearch.Port,
		User:            reg.Elasticsearch.Auth.Username,
		Pw:              reg.Elasticsearch.Auth.Password,
		DB:              reg.Elasticsearch.Database,
		InsecureSkipTLS: reg.Elasticsearch.TLS != nil && reg.Elasticsearch.TLS.InsecureSkipTLSVerify,
	}

	return Init(opts)
}

func Get() *elasticsearch.TypedClient {
	m.Lock()
	defer m.Unlock()

	return esClient
}
