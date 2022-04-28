package redis_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/baez90/nurse/config"
)

func PrepareRedisContainer(tb testing.TB) *config.Server {
	tb.Helper()

	const redisPort = "6379/tcp"

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	tb.Cleanup(cancel)

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "docker.io/redis:alpine",
			Name:         tb.Name(),
			ExposedPorts: []string{redisPort},
			WaitingFor:   wait.ForListeningPort(redisPort),
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

	u, err := url.Parse(fmt.Sprintf("%s/0?MaxRetries=3", ep))
	if err != nil {
		tb.Fatalf("url.Parse() err = %v", err)
	}

	srv, err := config.ParseFromURL(u)
	if err != nil {
		tb.Fatalf("config.ParseFromURL() err = %v", err)
	}
	return srv
}
