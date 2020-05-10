package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSendNotifyAboutTimetableUpdate(t *testing.T) {
	db.PrepareTestDatabase()

	var c int

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 8046))
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	router.HandleFunc(
		"/messages.send",
		func(w http.ResponseWriter, r *http.Request) {
			c += 1
			fmt.Fprint(
				w,
				`{"response": [{"peer_id": 1, "message_id": 1}, 
								 {"peer_id": 2, "message_id": 2, "error": {"code": 1, "description": "a"}}]}`,
			)
		},
	)

	ts := httptest.NewUnstartedServer(router)
	ts.Listener.Close()
	ts.Listener = l
	ts.Start()
	defer ts.Close()

	tools.Fatal(t, assert.NoError(t, db.CreateGroupIfNotExists("58-1")))
	group, err := db.GetGroupByGroupCode("58-1")
	tools.Fatal(t, assert.NoError(t, err))

	for i := 500; i < 864; i++ {
		if i < 50 {
			tools.Fatal(t, assert.NoError(t, db.CreateUser(i, group.Id)))
			user, err := db.GetUserByVkID(i)
			tools.Fatal(t, assert.NoError(t, err))

			if i%2 == 0 {
				tools.Fatal(t, assert.NoError(t, user.SetSubscribe(false)))
			} else {
				tools.Fatal(t, assert.NoError(t, user.SetActive(false)))
			}
		} else {
			tools.Fatal(t, assert.NoError(t, db.CreateUser(i, group.Id)))
		}
	}

	b, err := vk.CreateBot(vk.Settings{
		GroupID: 1,
		Token:   "s",
		BaseURL: ts.URL,
	})
	tools.Fatal(t, assert.NoError(t, err))

	sleepDuration = time.Nanosecond
	err = SendNotifyAboutTimetableUpdate(b)
	tools.Fatal(t, assert.NoError(t, err))
	assert.Equal(t, 4, c)
}
