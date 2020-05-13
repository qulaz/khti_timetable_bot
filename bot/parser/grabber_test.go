package parser

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"gopkg.in/khaiql/dbcleaner.v2"
	"gopkg.in/khaiql/dbcleaner.v2/engine"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type GrabberTestSuite struct {
	suite.Suite
}

var (
	Cleaner = dbcleaner.New()
	ts      *httptest.Server
)

func SetupTestServer() {
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

	ts = httptest.NewUnstartedServer(router)
	ts.Listener.Close()
	ts.Listener = l
	ts.Start()
}

func (suite *GrabberTestSuite) SetupSuite() {
	common.TestInits()
	db.InitTestDatabase()
	db := engine.NewPostgresEngine(helpers.Config.POSTGRES_DSN)
	Cleaner.SetEngine(db)
	SetupTestServer()
}

func (suite *GrabberTestSuite) TearDownSuite() {
	Cleaner.Clean("groups", "timetable", "users")
	db.CloseDatabase()
	ts.Close()
}

func (suite *GrabberTestSuite) SetupTest() {
	Cleaner.Acquire("groups", "timetable", "users")
	db.PrepareTestDatabase()
}

func (suite *GrabberTestSuite) TestGrabTimetableFromSite() {
	baseUrl = ts.URL

	filePath, err := GrabTimetableFromSite()
	tools.Fatal(suite.T(), suite.NoError(err))
	tools.Fatal(suite.T(), suite.Equal(filePath, "timetable.xls"))
	suite.NoError(os.Remove(filePath))
}

func (suite *GrabberTestSuite) TestUpdateTimetable_empty_db() {
	Cleaner.Clean("groups", "timetable", "users")
	Cleaner.Acquire("groups", "timetable", "users")

	baseUrl = ts.URL

	tools.Fatal(suite.T(), suite.NoError(UpdateTimetable(nil)))

	timetable, err := db.GetTimetable()
	tools.Fatal(suite.T(), suite.NoError(err))
	expectedTimetable, err := Parse("parser/testdata/different_timetable.xls")
	tools.Fatal(suite.T(), suite.NoError(err))
	suite.Equal(expectedTimetable, timetable)
	suite.NoError(os.Remove("timetable.xls"))
}

func (suite *GrabberTestSuite) TestUpdateTimetable_old_timetable() {
	baseUrl = ts.URL

	b, err := vk.CreateBot(vk.Settings{
		GroupID: 1,
		Token:   "token",
		BaseURL: ts.URL,
	})
	tools.Fatal(suite.T(), suite.NoError(err))
	tools.Fatal(suite.T(), suite.NoError(UpdateTimetable(b)))

	timetable, err := db.GetTimetable()
	tools.Fatal(suite.T(), suite.NoError(err))
	expectedTimetable, err := Parse("parser/testdata/different_timetable.xls")
	tools.Fatal(suite.T(), suite.NoError(err))
	suite.Equal(expectedTimetable, timetable)
	suite.NoError(os.Remove("timetable.xls"))
}

func TestGrabberSuite(t *testing.T) {
	suite.Run(t, new(GrabberTestSuite))
}
