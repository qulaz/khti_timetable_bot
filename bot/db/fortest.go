package db

import (
	"database/sql"
	"fmt"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/ory/dockertest/v3"
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
	"log"
	"os"
	"testing"
)

var (
	testDb   *sql.DB
	fixtures *testfixtures.Loader
)

// Функция TestMain для тестов, в которых используется база данных. В самих тестах необходимо вызвать функцию
// PrepareTestDatabase() из этого пакета для записи фикстур в тестовую базу
func TestMainWithDb(m *testing.M) {
	var pool *dockertest.Pool
	var resource *dockertest.Resource
	var dsn string

	// Инициализация конфигурации приложения
	common.TestInits()
	// Если тест запущен не в CI системе - создаем контейнер с Postgres и инициализируем подключение к созданной базе
	if !helpers.Config.IS_CI_TEST {
		pool, resource, dsn = common.CreatePostgres()
		testDb = InitDatabaseForTest(dsn)
	} else {
		// В системе CI необходимо указать переменные среды POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_HOST,
		// POSTGRES_PORT, POSTGRES_DB для корректной работы
		testDb = InitDatabaseForTest(
			fmt.Sprintf(
				"postgres://%s:%s@%s:%d/%s?sslmode=disable",
				helpers.Config.POSTGRES_USER,
				helpers.Config.POSTGRES_PASSWORD,
				helpers.Config.POSTGRES_HOST,
				helpers.Config.POSTGRES_PORT,
				helpers.Config.POSTGRES_DB,
			),
		)
	}

	var err error
	// Инициализация фикстур базы данных. В названии базы данных должно быть слово test для корректной инициализации
	fixtures, err = testfixtures.New(
		testfixtures.Database(testDb),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("db/testdata/fixtures"),
	)
	if err != nil {
		log.Println("Ошибка инициализации фикстур", err)
	}

	var code = 1
	if err == nil {
		code = m.Run()
	}

	if !helpers.Config.IS_CI_TEST {
		// Удаление тестового контейнера Postgres
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}

	os.Exit(code)
}

// Применение фикстур БД
func PrepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		log.Println("Ошибка применения фикстур", err)
	}
}
