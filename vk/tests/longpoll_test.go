package tests

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"gitlab.com/qulaz/khti_timetable_bot/vk/tests/consts"
	"gitlab.com/qulaz/khti_timetable_bot/vk/tests/handlers"
	"gitlab.com/qulaz/khti_timetable_bot/vk/tests/mocks"
	"gitlab.com/qulaz/khti_timetable_bot/vk/tools"
	"os"
	"testing"
	"time"
)

func init() {
	err := os.Chdir(consts.GetTestDir())
	if err != nil {
		panic(err)
	}
}

func TestLongPoll_GetUpdates(t *testing.T) {
	server := handlers.SetupTestServer()
	defer server.Close()
	consts.B.SetBaseUrl(server.URL)
	consts.ErrB.SetBaseUrl(server.URL)

	mocks.ResetMocks()
	tools.Now = func() time.Time { return mocks.TestTime }

	s, err := consts.B.GetLongPollServer(consts.TEST_LONG_POLL_WAIT)
	if err != nil {
		t.Fatalf("не удалось получить long poll сервер для теста GetUpdates: %q", err)
	}

	for n, testUpdatePack := range mocks.LongPollTestUpdates {
		updates, err := vk.GetUpdates(consts.B, s)
		if err != nil {
			t.Fatalf("%d: не удалось получить обновления: %q", n+1, err)
		}

		assert.Equal(t, len(testUpdatePack), len(updates))

		for i, update := range updates {
			expectedUpdate := testUpdatePack[i]
			assert.Equal(t, expectedUpdate, update)
		}
	}

	var longPollKeyNum = consts.CurrentLongPollKey
	if _, err := vk.GetUpdates(consts.B, s); err != nil {
		t.Fatalf("не удалось получить обновления: %q", err)
	}
	if longPollKeyNum == consts.CurrentLongPollKey {
		t.Error("При ошибке функция не меняет Long Poll сервер")
	}
}
