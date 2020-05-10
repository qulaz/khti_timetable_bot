package parser

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
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	db.TestMainWithDb(m)
}

func SetupTestServer() *httptest.Server {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "127.0.0.1", 8044))
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	router.HandleFunc("/obuchenie/raspisanie-zanyatiy.php", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "parser/testdata/khti_timetable_page.html")
	})
	router.HandleFunc(
		"/documents/Расписание/Расписание занятий, весенний семестр, 2019-2020 уч.гг. (ОФО).xls",
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "parser/testdata/different_timetable.xls")
		},
	)
	router.HandleFunc(
		"/messages.send",
		func(w http.ResponseWriter, r *http.Request) {
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

	return ts
}

func TestGrabTimetableFromSite(t *testing.T) {
	ts := SetupTestServer()
	defer ts.Close()

	baseUrl = ts.URL

	filePath, err := GrabTimetableFromSite()
	tools.Fatal(t, assert.NoError(t, err))
	tools.Fatal(t, assert.Equal(t, filePath, "timetable.xls"))
	assert.NoError(t, os.Remove(filePath))
}

func TestUpdateTimetable_empty_db(t *testing.T) {
	ts := SetupTestServer()
	defer ts.Close()
	baseUrl = ts.URL

	tools.Fatal(t, assert.NoError(t, UpdateTimetable(nil)))

	timetable, err := db.GetTimetable()
	tools.Fatal(t, assert.NoError(t, err))
	expectedTimetable, err := Parse("parser/testdata/different_timetable.xls")
	tools.Fatal(t, assert.NoError(t, err))
	assert.Equal(t, expectedTimetable, timetable)
	assert.NoError(t, os.Remove("timetable.xls"))
}

func TestUpdateTimetable_old_timetable(t *testing.T) {
	db.PrepareTestDatabase()
	ts := SetupTestServer()
	defer ts.Close()
	baseUrl = ts.URL

	b, err := vk.CreateBot(vk.Settings{
		GroupID: 1,
		Token:   "token",
		BaseURL: ts.URL,
	})
	tools.Fatal(t, assert.NoError(t, err))
	tools.Fatal(t, assert.NoError(t, UpdateTimetable(b)))

	timetable, err := db.GetTimetable()
	tools.Fatal(t, assert.NoError(t, err))
	expectedTimetable, err := Parse("parser/testdata/different_timetable.xls")
	tools.Fatal(t, assert.NoError(t, err))
	assert.Equal(t, expectedTimetable, timetable)
	assert.NoError(t, os.Remove("timetable.xls"))
}
