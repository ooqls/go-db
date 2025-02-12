package testutils

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/braumsmilk/go-db"
	"github.com/braumsmilk/go-db/init/seed"
	"github.com/braumsmilk/go-registry"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var reg registry.Registry = registry.Registry{}

func InitRedis() testcontainers.Container {
	ctx := context.Background()
	c := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379"},
		WaitingFor:   &wait.LogStrategy{Log: "Ready to accept connections"},
		Env: map[string]string{
			"REDIS_PASSWORD": "password",
		},
	}

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

	reg.Redis = &registry.Server{
		Host: "localhost",
		Port: port.Int(),
		Auth: registry.Auth{
			Enabled:  true,
			Password: "password",
		},
	}
	registry.Set(reg)

	return container
}

func InitPostgres(tableStmts []string, indexStmts []string) testcontainers.Container {
	ctx := context.Background()
	arch := runtime.GOARCH
	image := "postgres:latest"
	if arch == "arm64" {
		image = "arm64v8/postgres:latest"
	}
	log.Printf("Detected architecture: %s", arch)
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

	gc := testcontainers.GenericContainerRequest{
		ContainerRequest: c,
		Started:          true,
	}
	container, err := testcontainers.GenericContainer(ctx, gc)
	if err != nil {
		panic(fmt.Errorf("failed to start postgres container: %v", err))
	}
	// host, err := container.Host(ctx)
	// if err != nil {
	// 	panic(fmt.Errorf("failed to get host ip from container: %v", err))
	// }

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		panic(fmt.Errorf("failed to get mapped postgres port: %v", err))
	}

	log.Printf("postgres should be running at localhost:%d", port.Int())
	time.Sleep(time.Second * 5)

	
	reg.Postgres = &registry.Server{
		Host: "localhost",
		Port: port.Int(),
		Auth: registry.Auth{
			Enabled:  true,
			Username: "user",
			Password: "user100",
		},
	}
	registry.Set(reg)

	err = db.InitDefault()
	if err != nil {
		panic(err)
	}

	seed.SeedPostgresDatabase(tableStmts, indexStmts)

	return container
}
