package helpers

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
	"path"
	"runtime"
)

var Config *config = &config{}

type config struct {
	POSTGRES_HOST     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	POSTGRES_PORT     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	POSTGRES_USER     string `env:"POSTGRES_USER" envDefault:"postgres"`
	POSTGRES_DB       string `env:"POSTGRES_DB" envDefault:"postgres"`
	POSTGRES_PASSWORD string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
	POSTGRES_DSN      string `env:"POSTGRES_DSN"`

	IS_DEBUG    bool   `env:"DEBUG" envDefault:"false"`
	LOG_LEVEL   string `env:"LOG_LEVEL" envDefault:"INFO"`
	PROJECT_DIR string

	SENTRY_DSN string `env:"SENTRY_DSN"`

	VK_GROUP_TOKEN string `env:"VK_GROUP_TOKEN"`
	VK_GROUP_ID    int    `env:"VK_GROUP_ID"`
	VK_ADMIN_ID    int    `env:"VK_ADMIN_ID"`
}

// Инициализация конфига
func LoadConfigFromEnv() {
	// Загрузка переменных из файла .env в переменные среды
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Не найден .env файл. Заполните .dist.env файл и переименуйте его в .env:")
	}

	// Перенос переменных среды на Go-структуру конфига
	if err := env.Parse(Config); err != nil {
		log.Fatal("Ошибка маппинга переменных среды: ", err)
	}

	// Установка директории проекта
	Config.PROJECT_DIR = GetProjectRootDir()
}

// Возвращает полный путь до директории bot
func GetProjectRootDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(filename), "..")
}
