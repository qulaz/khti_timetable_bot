package db

import (
	"database/sql"
	"github.com/go-testfixtures/testfixtures/v3"
	"log"
)

var (
	testDb   *sql.DB
	fixtures *testfixtures.Loader
)

func InitTestDatabase() {
	testDb = InitDatabase()

	var err error
	// Инициализация фикстур базы данных. В названии базы данных должно быть слово test для корректной инициализации
	fixtures, err = testfixtures.New(
		testfixtures.Database(testDb),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("db/testdata/fixtures"),
	)
	if err != nil {
		log.Fatal("Ошибка инициализации фикстур", err)
	}
}

// Применение фикстур БД
func PrepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		log.Println("Ошибка применения фикстур", err)
	}
}
