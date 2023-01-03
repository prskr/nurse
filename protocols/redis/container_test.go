package redis_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"code.icb4dc0.de/prskr/nurse/config"
)

func PrepareRedisContainer(tb testing.TB) *config.Server {
	tb.Helper()
	const redisPort = "6379/tcp"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	tb.Cleanup(cancel)

	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Image:        "docker.io/redis:alpine",
			ExposedPorts: []string{redisPort},
			SkipReaper:   true,
			AutoRemove:   true,
			WaitingFor:   wait.ForListeningPort(redisPort),
		},
		Started: true,
		Logger:  tc.TestLogger(tb),
	})
	if err != nil {
		tb.Fatalf("tc.GenericContainer() err = %v", err)
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
