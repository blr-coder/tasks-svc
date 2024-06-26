package psql_store

import (
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

var dbConnTest *sqlx.DB

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		return
	}

	// Run docker container for storage testing
	dockerForTest := MustDockerForStorageTest()

	//
	MustRunPsqlStorageTest(m, dockerForTest, "../psql_store/migrations", func(conn *sqlx.DB) {
		dbConnTest = conn

		cleanup()
	})
}

func cleanup() {
	dbConnTest.MustExec("DELETE FROM tasks")
}

func MustRunPsqlStorageTest(m *testing.M, dockerForTest *DockerForStorageTests, pathToMigrations string, prepareFunc func(conn *sqlx.DB)) {
	_, filename, _, _ := runtime.Caller(0)
	migrationDir, err := filepath.Abs(filepath.Join(path.Dir(filename), pathToMigrations))
	if err != nil {
		log.Fatalf("could not get migrations dir: %s", err.Error())
	}

	//dockerForTest := MustDockerForStorageTest()

	conn, err := dockerForTest.GetPostgresConn(migrationDir)
	if err != nil {
		log.Fatalf("failed to run postgres container: %s", err.Error())
	}

	prepareFunc(conn)

	code := m.Run()
	dockerForTest.Purge()
	os.Exit(code)
}

type DockerForStorageTests struct {
	pool *dockertest.Pool

	postgres *postgresDB
}

func (dt *DockerForStorageTests) GetPostgresConn(migrationsDir string) (*sqlx.DB, error) {
	dt.postgres = &postgresDB{
		migrationsDir: migrationsDir,
	}

	if err := dt.postgres.runAndUpMigrations(dt.pool); err != nil {
		return nil, err
	}
	return dt.postgres.conn, nil
}

func (dt *DockerForStorageTests) Purge() error {
	if dt.postgres != nil {
		if err := dt.pool.Purge(dt.postgres.resource); err != nil {
			return fmt.Errorf("could not purge postgres resource: %s", err)
		}
	}
	return nil
}

func MustDockerForStorageTest() *DockerForStorageTests {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker pool: %s", err)
	}

	return &DockerForStorageTests{
		pool: pool,
	}
}

type postgresDB struct {
	resource *dockertest.Resource

	conn          *sqlx.DB
	migrationsDir string
	databaseURL   string
}

func (p *postgresDB) runAndUpMigrations(pool *dockertest.Pool) error {
	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13-alpine",
		Env: []string{
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=postgres",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return fmt.Errorf("could not start resource: %w", err)
	}

	p.resource = resource

	hostAndPort := resource.GetHostPort("5432/tcp")
	if os.Getenv("GITLAB_CI") != "" {
		hostAndPort = fmt.Sprintf("%s:%s", "172.17.0.1", resource.GetPort("5432/tcp"))
	}

	databaseURL := fmt.Sprintf("postgres://postgres:postgres@%s/postgres?sslmode=disable", hostAndPort)

	log.Println("connecting to database on url: ", databaseURL)

	p.databaseURL = databaseURL

	const wait = 120

	_ = resource.Expire(wait)
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = wait * time.Second

	if err := pool.Retry(func() error {
		p.conn, err = sqlx.Open("postgres", databaseURL)
		if err != nil {
			return err
		}
		return p.conn.Ping()
	}); err != nil {
		return fmt.Errorf("could not connect to docker: %w", err)
	}

	return p.migrateUp()
}

/*func (p *postgres) migrate() error {
	m, err := migrate.NewWithSourceInstance(
		p.migrationsDir,
		p.conn,
	)
	if err != nil {
		return fmt.Errorf("could not apply migrations: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error when migration up: %v", err)
	}

	log.Println("migration completed!")

	return nil
}*/

func (p *postgresDB) migrateUp() error {
	// Настройка драйвера PostgreSQL
	driver, err := postgres.WithInstance(p.conn.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %w", err)
	}

	// Настройка источника миграций из директории
	src, err := (&file.File{}).Open("file://" + p.migrationsDir)
	if err != nil {
		return fmt.Errorf("could not open file source: %w", err)
	}

	// Создание экземпляра migrate
	m, err := migrate.NewWithInstance("file", src, "postgres", driver)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	// Применение миграций
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not apply migrations: %w", err)
	}

	return nil
}
