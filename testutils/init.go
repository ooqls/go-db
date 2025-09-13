package testutils

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ooqls/go-registry"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var reg registry.Registry = registry.Registry{}

const (
	opt_logging = "logging"
	opt_env     = "env"
	opt_tag     = "tags"
)

type options struct {
	key   string
	value interface{}
}

func WithLogging() options {
	return options{
		key:   opt_logging,
		value: true,
	}
}

func WithEnv(key string, value map[string]string) options {
	return options{
		key:   key,
		value: value,
	}
}

func WithTags(tag string) options {
	return options{
		key:   opt_tag,
		value: tag,
	}
}

func isArm64() bool {
	arch := runtime.GOARCH
	return arch == "arm64"
}

func applyOptions(c *testcontainers.ContainerRequest, opts ...options) {
	for _, opt := range opts {
		switch opt.key {
		case opt_logging:
			w := c.BuildLogWriter()
			io.Copy(w, os.Stdout)
		case opt_env:
			envMap := opt.value.(map[string]string)
			maps.Copy(c.Env, envMap)
		case opt_tag:
			baseImage := strings.Split(c.Image, ":")[0]
			c.Image = fmt.Sprintf("%s:%s", baseImage, opt.value)
		}
	}
}

func StartRedis(ctx context.Context, opts ...options) testcontainers.Container {
	c := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379"},
		WaitingFor:   &wait.LogStrategy{Log: "Ready to accept connections"},
		Env: map[string]string{
			"REDIS_PASSWORD": "password",
		},
	}

	applyOptions(&c, opts...)

	gc := testcontainers.GenericContainerRequest{
		ContainerRequest: c,
		Started:          true,
	}

	container, err := testcontainers.GenericContainer(ctx, gc)
	if err != nil {
		panic(fmt.Errorf("failed to start redis container: %v", err))
	}

	port, err := container.MappedPort(ctx, "6379")
	if err != nil {
		panic(fmt.Errorf("failed to get mapped redis port: %v", err))
	}

	log.Printf("redis should be running at localhost:%d", port.Int())
	time.Sleep(time.Second * 5)

	reg.Redis = &registry.Database{
		Database: "0",
		Server: registry.Server{
			Host: "localhost",
			Port: port.Int(),

			Auth: registry.Auth{
				Enabled:  true,
				Password: "password",
			},
		},
	}
	registry.Set(reg)

	return container
}

func StartPostgres(ctx context.Context, opts ...options) testcontainers.Container {
	image := "postgres:latest"
	if isArm64() {
		image = "arm64v8/postgres:latest"
	}

	c := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{"5432"},
		Env: map[string]string{
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "user100",
			"POSTGRES_DB":       "postgres",
		},
		WaitingFor: &wait.LogStrategy{Log: "database system is ready to accept connections"},
	}

	applyOptions(&c, opts...)

	gc := testcontainers.GenericContainerRequest{
		ContainerRequest: c,
		Started:          true,
	}
	container, err := testcontainers.GenericContainer(ctx, gc)
	if err != nil {
		panic(fmt.Errorf("failed to start postgres container: %v", err))
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		panic(fmt.Errorf("failed to get mapped postgres port: %v", err))
	}

	log.Printf("postgres should be running at localhost:%d", port.Int())
	time.Sleep(time.Second * 5)

	reg.Postgres = &registry.Database{
		Database: "test",
		Server: registry.Server{
			Host: "localhost",
			Port: port.Int(),
			Auth: registry.Auth{
				Enabled:  true,
				Username: "user",
				Password: "user100",
			},
		},
	}
	registry.Set(reg)

	return container
}

func StartElasticsearch(ctx context.Context, opts ...options) testcontainers.Container {
	image := "elasticsearch:8.18.0"
	if isArm64() {
		image = "arm64v8/elasticsearch:8.18.0"
	}

	c := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{"9200"},
		Env: map[string]string{
			"ELASTIC_PASSWORD": "changeme",
			"discovery.type":   "single-node",
			"ES_JAVA_OPTS":     "-Xms512m -Xmx512m",
		},
		WaitingFor: wait.ForHTTP("/_cluster/health").
			WithBasicAuth("elastic", "changeme").
			WithAllowInsecure(true).
			WithMethod("GET").
			WithTLS(true, &tls.Config{InsecureSkipVerify: true}).
			WithStatusCodeMatcher(func(status int) bool {
				return status == http.StatusOK
			}),
	}

	applyOptions(&c, opts...)

	gc := testcontainers.GenericContainerRequest{
		ContainerRequest: c,
		Started:          true,
	}

	container, err := testcontainers.GenericContainer(ctx, gc)
	if err != nil {
		panic(fmt.Errorf("failed to start elasticsearch container: %v", err))
	}

	time.Sleep(time.Second * 10)

	port, err := container.MappedPort(ctx, "9200")
	if err != nil {
		panic(fmt.Errorf("failed to get mapped elasticsearch port: %v", err))
	}

	log.Printf("elasticsearch should be running at localhost:%d", port.Int())

	reg.Elasticsearch = &registry.Database{
		Database: "elasticsearch",
		Server: registry.Server{
			Host: "localhost",
			Port: port.Int(),
			Auth: registry.Auth{
				Enabled:  true,
				Password: "changeme",
				Username: "elastic",
			},
			TLS: &registry.TLSConfig{
				Enabled:               true,
				InsecureSkipTLSVerify: true,
			},
		},
	}
	registry.Set(reg)

	return container
}
