package common

import (
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
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
