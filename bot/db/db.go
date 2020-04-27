package db

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
	"log"
	"time"
)

var db *sql.DB

type scanner interface {
	Scan(dest ...interface{}) error
}

type queryable interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// Ожидание готовности базы данных принимать соединения
func waitDatabase(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	for {
		err = db.Ping()
		if err == nil {
			break
		}
		helpers.Logger.Info("Database is NOT ready!")
		time.Sleep(time.Second * 2)
	}
	helpers.Logger.Info("Database is ready!")

	return db
}

// Применение миграций
func applyMigrations() {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		helpers.Logger.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver,
	)
	if err != nil {
		helpers.Logger.Fatal(err)
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		helpers.Logger.Fatal(err)
	}
}

// Инициализация базы данных
func InitDatabase() *sql.DB {
	db = waitDatabase(helpers.Config.POSTGRES_DSN)
	applyMigrations()

	return db
}

// Инициализация базы данных для тестов
func InitDatabaseForTest(dsn string) *sql.DB {
	fmt.Println(dsn)
	db = waitDatabase(dsn)
	applyMigrations()

	return db
}

// Закрытие соединения с базой данных
func CloseDatabase() {
	db.Close()
}
