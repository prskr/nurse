package redis_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"

	"code.1533b4dc0.de/prskr/nurse/config"
)

func PrepareRedisContainer(tb testing.TB) *config.Server {
	tb.Helper()
	const redisPort = "6379/tcp"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	tb.Cleanup(cancel)

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "docker.io/redis:alpine",
			ExposedPorts: []string{redisPort},
			SkipReaper:   true,
			AutoRemove:   true,
		},
		Started: true,
		Logger:  testcontainers.TestLogger(tb),
	})
	if err != nil {
		tb.Fatalf("testcontainers.GenericContainer() err = %v", err)
	}

	tb.Cleanup(func() {
		if err := container.Terminate(context.Background()); err != nil {
			tb.Errorf("container.Terminate() err = %v", err)
		}
	})

	ep, err := container.PortEndpoint(ctx, redisPort, "redis")
	if err != nil {
		tb.Fatalf("container.PortEndpoint() err = %v", err)
	}

	srv := new(config.Server)
	if err := srv.UnmarshalURL(fmt.Sprintf("%s/0?MaxRetries=3", ep)); err != nil {
		tb.Fatalf("config.ParseFromURL() err = %v", err)
	}
	return srv
}
