package common

import (
	"database/sql"
	"fmt"
	"github.com/ory/dockertest/v3"
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
	"log"
	"os"
)

// Инициализация всей инфраструктуры приложения, кроме базы данных
func TestInits() {
	if err := os.Chdir(helpers.GetProjectRootDir()); err != nil {
		panic(err)
	}
	helpers.LoadConfigFromEnv()
	helpers.InitLogger()
}

// Создание контейнера Postgres для тестов
func CreatePostgres() (*dockertest.Pool, *dockertest.Resource, string) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "12-alpine", []string{
		"POSTGRES_PORT=5432",
		"POSTGRES_USER=test",
		"POSTGRES_DB=test",
		"POSTGRES_PASSWORD=test",
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	dsn := fmt.Sprintf("postgres://test:test@localhost:%s/test?sslmode=disable", resource.GetPort("5432/tcp"))

	if err := pool.Retry(func() error {
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return err
		}
		err = db.Ping()
		if err == nil {
			db.Close()
		}
		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return pool, resource, dsn
}
