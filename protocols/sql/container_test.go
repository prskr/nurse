package sql_test

import (
	"context"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/td"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"code.icb4dc0.de/prskr/nurse/config"
	"code.icb4dc0.de/prskr/nurse/internal/values"
)

const (
	dbName     = "nurse"
	dbUser     = "nurse"
	dbPassword = "asi1EeYi"
)

func PreparePostgresContainer(tb testing.TB) (name string, cfg *config.Server) {
	tb.Helper()

	const postgresPort = "5432/tcp"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	tb.Cleanup(cancel)

	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Image:        "docker.io/postgres:alpine",
			ExposedPorts: []string{postgresPort},
			SkipReaper:   true,
			AutoRemove:   true,
			Env: map[string]string{
				"POSTGRES_USER":     dbUser,
				"POSTGRES_PASSWORD": dbPassword,
				"POSTGRES_DB":       dbName,
			},
			WaitingFor: wait.ForListeningPort(postgresPort),
		},
		Started: true,
		Logger:  tc.TestLogger(tb),
	})

	td.CmpNoError(tb, err, "tc.GenericContainer()")

	tb.Cleanup(func() {
		td.CmpNoError(tb, container.Terminate(context.Background()), "container.Terminate()")
	})

	ep, err := container.PortEndpoint(ctx, postgresPort, "postgres")
	td.CmpNoError(tb, err, "container.PortEndpoint()")

	srv := new(config.Server)
	td.CmpNoError(tb, srv.UnmarshalURL(ep), "srv.UnmarshalURL()")

	srv.Path = append(srv.Path, dbName)
	srv.Credentials = &config.Credentials{
		Username: dbUser,
		Password: values.StringP(dbPassword),
	}

	name, err = container.Name(ctx)
	td.CmpNoError(tb, err, "container.Name()")

	return name, srv
}

func PrepareMariaDBContainer(tb testing.TB) (name string, cfg *config.Server) {
	tb.Helper()

	const mysqlPort = "3306/tcp"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	tb.Cleanup(cancel)

	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: tc.ContainerRequest{
			Image:        "docker.io/mariadb:10",
			ExposedPorts: []string{mysqlPort},
			SkipReaper:   true,
			AutoRemove:   true,
			Env: map[string]string{
				"MARIADB_USER":                 dbUser,
				"MARIADB_PASSWORD":             dbPassword,
				"MARIADB_RANDOM_ROOT_PASSWORD": "1",
				"MARIADB_DATABASE":             dbName,
			},
			WaitingFor: wait.ForListeningPort(mysqlPort),
		},
		Started: true,
		Logger:  tc.TestLogger(tb),
	})

	td.CmpNoError(tb, err, "tc.GenericContainer()")

	tb.Cleanup(func() {
		td.CmpNoError(tb, container.Terminate(context.Background()), "container.Terminate()")
	})

	ep, err := container.PortEndpoint(ctx, mysqlPort, "mysql")
	td.CmpNoError(tb, err, "container.PortEndpoint()")

	srv := new(config.Server)
	td.CmpNoError(tb, srv.UnmarshalURL(ep), "srv.UnmarshalURL()")

	srv.Path = append(srv.Path, dbName)
	srv.Credentials = &config.Credentials{
		Username: dbUser,
		Password: values.StringP(dbPassword),
	}

	name, err = container.Name(ctx)
	td.CmpNoError(tb, err, "container.Name()")

	return name, srv
}
