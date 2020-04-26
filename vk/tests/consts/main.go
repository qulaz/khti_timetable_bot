package consts

import (
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"path"
	"runtime"
)

const (
	TEST_GROUP_ID    = 52
	TEST_TOKEN       = "test_token"
	TEST_API_VERSION = vk.VK_API_VERISON
	// Все запросы с этой версией API будут выполняться с ошибкой
	TEST_ERR_API_VERSION = "8.81"
)

var B, _ = vk.CreateBot(vk.Settings{
	GroupID:    TEST_GROUP_ID,
	Token:      TEST_TOKEN,
	ApiVersion: TEST_API_VERSION,
})

// На все запросы с этого бота возвращаются ошибки
var ErrB, _ = vk.CreateBot(vk.Settings{
	GroupID:    TEST_GROUP_ID,
	Token:      TEST_TOKEN,
	ApiVersion: TEST_ERR_API_VERSION,
})

func GetTestDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(filename), "..")
}
