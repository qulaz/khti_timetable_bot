package handlers

import (
	"fmt"
	"gitlab.com/qulaz/khti_timetable_bot/vk/tests/consts"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
)

const (
	TEST_HOST = "127.0.0.1"
	TEST_PORT = 8045
)

func init() {
	err := os.Chdir(consts.GetTestDir())
	if err != nil {
		panic(err)
	}
}

func SetupTestServer() *httptest.Server {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", TEST_HOST, TEST_PORT))
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	router.HandleFunc("/groups.getLongPollServer", GetLongPollServerHandler)
	router.HandleFunc("/messages.send", MessageSendHandler)
	router.HandleFunc(consts.TEST_LONG_POLL_SERVER_PREFIX, LongPollServerHandler)
	router.HandleFunc("/some.method", SomeMethodHandler)

	ts := httptest.NewUnstartedServer(router)
	ts.Listener.Close()
	ts.Listener = l
	ts.Start()

	return ts
}
