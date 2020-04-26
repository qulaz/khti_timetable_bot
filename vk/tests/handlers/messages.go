package handlers

import (
	"fmt"
	"gitlab.com/qulaz/khti_timetable_bot/vk/tests/consts"
	"gitlab.com/qulaz/khti_timetable_bot/vk/tests/helpers"
	"log"
	"net/http"
	"os"
)

func init() {
	err := os.Chdir(consts.GetTestDir())
	if err != nil {
		panic(err)
	}
}

func MessageSendHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal()
	}

	token := r.Form.Get("access_token")
	apiVersion := r.Form.Get("v")
	userIDs := r.Form.Get("user_ids")

	if apiVersion == consts.TEST_ERR_API_VERSION {
		resp := helpers.LoadStringFromFile("api_test/messages.send/error.json")
		if _, err := fmt.Fprintln(w, resp); err != nil {
			log.Fatal(err)
		}
		return
	}

	if token != consts.TEST_TOKEN || apiVersion != consts.TEST_API_VERSION {
		log.Fatalf("Не совпадают токен или версия api. token=%q v=%q", token, apiVersion)
	}

	var resp string
	if userIDs != "" {
		// Ответ при отправке сообщения нескольким пользователям
		resp = helpers.LoadStringFromFile("api_test/messages.send/multiple_messages_responsef.json")
		resp = fmt.Sprintf(resp, consts.TestPeerID1, consts.TestPeerID2)
	} else {
		// Обычный ответ
		resp = fmt.Sprintf(`{"response": %d}`, consts.TestMessageID)
	}

	if _, err := fmt.Fprintln(w, resp); err != nil {
		log.Fatal(err)
	}
}
