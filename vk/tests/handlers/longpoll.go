package handlers

import (
	"fmt"
	"gitlab.com/qulaz/khti_timetable_bot/vk/tests/consts"
	"gitlab.com/qulaz/khti_timetable_bot/vk/tests/helpers"
	"log"
	"net/http"
	"os"
	"strconv"
)

func init() {
	err := os.Chdir(consts.GetTestDir())
	if err != nil {
		panic(err)
	}
}

func GetLongPollServerHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal()
	}

	token := r.Form.Get("access_token")
	groupID := r.Form.Get("group_id")
	apiVersion := r.Form.Get("v")

	if apiVersion == consts.TEST_ERR_API_VERSION {
		resp := helpers.LoadStringFromFile("api_test/groups.getLongPollServer/error.json")
		if _, err := fmt.Fprintln(w, resp); err != nil {
			log.Fatal(err)
		}
		return
	}

	if token != consts.TEST_TOKEN || groupID != strconv.Itoa(consts.TEST_GROUP_ID) || apiVersion != consts.TEST_API_VERSION {
		log.Fatalf(
			"Не совпадают токен, группа или версия api. token=%q group_id=%q v=%q", token, groupID, apiVersion,
		)
	}

	respf := `{"response": {"key": "%s", "server": "http://%s:%d%s", "ts": "1"}}`

	if consts.CurrentLongPollKey == 1 {
		if _, err := fmt.Fprintf(
			w,
			respf,
			consts.TEST_LONG_POLL_KEY2,
			TEST_HOST,
			TEST_PORT,
			consts.TEST_LONG_POLL_SERVER_PREFIX,
		); err != nil {
			log.Fatal(err)
		}
		consts.CurrentLongPollKey = 2
	} else if consts.CurrentLongPollKey == 2 {
		if _, err := fmt.Fprintf(
			w,
			respf,
			consts.TEST_LONG_POLL_KEY1,
			TEST_HOST,
			TEST_PORT,
			consts.TEST_LONG_POLL_SERVER_PREFIX,
		); err != nil {
			log.Fatal(err)
		}
		consts.CurrentLongPollKey = 1
	}
}

func LongPollServerHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	wait := r.URL.Query().Get("wait")
	act := r.URL.Query().Get("act")
	ts := r.URL.Query().Get("ts")

	if (key != consts.TEST_LONG_POLL_KEY1 && key != consts.TEST_LONG_POLL_KEY2) || wait != strconv.Itoa(consts.TEST_LONG_POLL_WAIT) || act != "a_check" {
		log.Fatalf(
			"Не совпадают key, wait или act. key=%q wait=%q act=%q",
			key, wait, act,
		)
	}

	var resp string
	switch ts {
	case "1":
		resp = helpers.LoadStringFromFile("longpoll_test/ts1.json")
	case "3":
		resp = helpers.LoadStringFromFile("longpoll_test/ts3.json")
	case "6":
		resp = `{"failed":1, "ts":6}`
	}

	if _, err := fmt.Fprintln(w, resp); err != nil {
		log.Fatal(err)
	}
}
