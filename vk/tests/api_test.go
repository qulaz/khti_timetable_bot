package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"gitlab.com/qulaz/khti_timetable_bot/vk/tests/consts"
	"gitlab.com/qulaz/khti_timetable_bot/vk/tests/handlers"
	"os"
	"sort"
	"strings"
	"testing"
)

func init() {
	err := os.Chdir(consts.GetTestDir())
	if err != nil {
		panic(err)
	}
}

func TestAPI_RawRequest(t *testing.T) {
	server := handlers.SetupTestServer()
	defer server.Close()
	consts.B.SetBaseUrl(server.URL)

	params := vk.H{"some_key": "some_value"}
	raw, err := consts.B.RawRequest("some.method", params)
	if err != nil {
		t.Fatal("Ошибка при запросе:", err)
	}

	params["access_token"] = consts.TEST_TOKEN
	params["v"] = consts.TEST_API_VERSION

	// Sorting params keys
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	expectedString := ""
	for _, key := range keys {
		expectedString += fmt.Sprintf("%s=%q;", key, params[key])
	}

	assert.Equal(t, expectedString, strings.Trim(string(raw), "\n"))
}

func TestAPI_getLongPollServer(t *testing.T) {
	server := handlers.SetupTestServer()
	defer server.Close()
	consts.B.SetBaseUrl(server.URL)
	consts.ErrB.SetBaseUrl(server.URL)

	s, err := consts.B.GetLongPollServer(consts.TEST_LONG_POLL_WAIT)
	if err != nil {
		t.Fatal("Ошибка получения сервера:", err)
	}

	assert.Equal(t, server.URL+consts.TEST_LONG_POLL_SERVER_PREFIX, s.Server)
	assert.Equal(t, consts.TEST_LONG_POLL_WAIT, s.Wait)
	if s.Key != consts.TEST_LONG_POLL_KEY1 && s.Key != consts.TEST_LONG_POLL_KEY2 {
		t.Errorf(
			"key тестового long poll сервера не совподает с ожидаемым. тест=%q ожидание=%q/%q",
			s.Key,
			consts.TEST_LONG_POLL_KEY1,
			consts.TEST_LONG_POLL_KEY2,
		)
	}

	_, err = consts.ErrB.GetLongPollServer(consts.TEST_LONG_POLL_WAIT)
	if err == nil {
		t.Error("Не вернулась ожидаемая ошибка при запросе Long Poll сервера c `бракованного` бота")
	}
}

func TestAPI_SendMessage(t *testing.T) {
	server := handlers.SetupTestServer()
	defer server.Close()
	consts.B.SetBaseUrl(server.URL)
	consts.ErrB.SetBaseUrl(server.URL)

	opt := &vk.SendOptions{
		PeerID:  151,
		Message: "ыы",
	}

	mid, err := consts.B.SendMessage(opt)
	if err != nil {
		t.Fatal("Ошибка при отправке сообщения:", err)
	}

	assert.Equal(t, consts.TestMessageID, mid)

	_, err = consts.ErrB.SendMessage(opt)
	if err == nil {
		t.Error("Не вернулась ожидаемая ошибка при отправке сообщения от `бракованного` бота")
	}
}

func TestAPI_SendMultipleMessages(t *testing.T) {
	server := handlers.SetupTestServer()
	defer server.Close()
	consts.B.SetBaseUrl(server.URL)

	consts.B.SetBaseUrl(server.URL)
	consts.ErrB.SetBaseUrl(server.URL)

	opt := &vk.SendOptions{
		Message:        "ыы",
		IsSendMultiple: true,
	}

	resp, err := consts.B.SendMultipleMessages(opt, consts.TestPeerID1, consts.TestPeerID2)
	if err != nil {
		t.Fatal("Ошибка при отправке сообщения:", err)
	}
	for _, m := range resp {
		if m.PeerID != consts.TestPeerID1 && m.PeerID != consts.TestPeerID2 {
			t.Errorf(
				"Peer отправленного сообщения не совподает с ожидаемым. тест=%q ожидание=%q или %q",
				m.PeerID,
				consts.TestPeerID1,
				consts.TestPeerID2,
			)
		}
	}

	_, err = consts.ErrB.SendMultipleMessages(opt, consts.TestPeerID1, consts.TestPeerID2)
	if err == nil {
		t.Error("Не вернулась ожидаемая ошибка при отправке сообщения от `бракованного` бота")
	}
}

func TestAPI_SendTextMessage(t *testing.T) {
	server := handlers.SetupTestServer()
	defer server.Close()
	consts.B.SetBaseUrl(server.URL)

	consts.B.SetBaseUrl(server.URL)
	consts.ErrB.SetBaseUrl(server.URL)

	mid, err := consts.B.SendTextMessage("ыы", consts.TestPeerID1)
	if err != nil {
		t.Fatal("Ошибка при отправке сообщения:", err)
	}

	assert.Equal(t, consts.TestMessageID, mid)

	_, err = consts.ErrB.SendTextMessage("ыы", consts.TestPeerID1)
	if err == nil {
		t.Error("Не вернулась ожидаемая ошибка при отправке сообщения от `бракованного` бота")
	}
}

func TestAPI_SendMultipleTextMessages(t *testing.T) {
	server := handlers.SetupTestServer()
	defer server.Close()
	consts.B.SetBaseUrl(server.URL)

	consts.B.SetBaseUrl(server.URL)
	consts.ErrB.SetBaseUrl(server.URL)

	resp, err := consts.B.SendMultipleTextMessages("ыы", consts.TestPeerID1, consts.TestPeerID2)
	if err != nil {
		t.Fatal("Ошибка при отправке сообщения:", err)
	}
	for _, m := range resp {
		if m.PeerID != consts.TestPeerID1 && m.PeerID != consts.TestPeerID2 {
			t.Errorf(
				"Peer отправленного сообщения не совподает с ожидаемым. тест=%q ожидание=%q или %q",
				m.PeerID,
				consts.TestPeerID1,
				consts.TestPeerID2,
			)
		}
	}

	_, err = consts.ErrB.SendMultipleTextMessages("ыы", consts.TestPeerID1, consts.TestPeerID2)
	if err == nil {
		t.Error("Не вернулась ожидаемая ошибка при отправке сообщения от `бракованного` бота")
	}
}

func TestAPI_SendKeyboardMessage(t *testing.T) {
	server := handlers.SetupTestServer()
	defer server.Close()
	consts.B.SetBaseUrl(server.URL)

	consts.B.SetBaseUrl(server.URL)
	consts.ErrB.SetBaseUrl(server.URL)

	k := vk.NewKeyboard(false)
	k.AddTextButton("", "", nil)
	k.AddTextButton("", "", nil)
	k.AddTextButton("", "", nil)

	mid, err := consts.B.SendKeyboardMessage("ыы", k, consts.TestPeerID1)
	if err != nil {
		t.Fatal("Ошибка при отправке сообщения:", err)
	}

	assert.Equal(t, consts.TestMessageID, mid)

	_, err = consts.ErrB.SendKeyboardMessage("ыы", k, consts.TestPeerID1)
	if err == nil {
		t.Error("Не вернулась ожидаемая ошибка при отправке сообщения от `бракованного` бота")
	}
}

func TestAPI_SendMultipleKeyboardMessages(t *testing.T) {
	server := handlers.SetupTestServer()
	defer server.Close()
	consts.B.SetBaseUrl(server.URL)

	consts.B.SetBaseUrl(server.URL)
	consts.ErrB.SetBaseUrl(server.URL)

	k := vk.NewKeyboard(false)
	k.AddTextButton("", "", nil)
	k.AddTextButton("", "", nil)
	k.AddTextButton("", "", nil)

	resp, err := consts.B.SendMultipleKeyboardMessages("ыы", k, consts.TestPeerID1, consts.TestPeerID2)
	if err != nil {
		t.Fatal("Ошибка при отправке сообщения:", err)
	}
	for _, m := range resp {
		if m.PeerID != consts.TestPeerID1 && m.PeerID != consts.TestPeerID2 {
			t.Errorf(
				"Peer отправленного сообщения не совподает с ожидаемым. тест=%q ожидание=%q или %q",
				m.PeerID,
				consts.TestPeerID1,
				consts.TestPeerID2,
			)
		}
	}

	_, err = consts.ErrB.SendMultipleKeyboardMessages("ыы", k, consts.TestPeerID1, consts.TestPeerID2)
	if err == nil {
		t.Error("Не вернулась ожидаемая ошибка при отправке сообщения от `бракованного` бота")
	}
}
