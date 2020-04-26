package vk

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/qulaz/khti_timetable_bot/vk/tools"
	"io/ioutil"
	"net/url"
	"reflect"
	"strings"
)

const (
	VK_API_VERISON  = "5.103"
	VK_API_BASE_URL = "https://api.vk.com/method/"
)

// Сокращение для типа map[string]interface{}. Используется, в основном, для написания параметов к POST-запросам
type H map[string]interface{}

type ApiError struct {
	ErrorCode     int                 `json:"error_code"`
	ErrorMsg      string              `json:"error_msg"`
	RequestParams []map[string]string `json:"request_params"`
}

// Приведение параметров разных типов к единому url.Values типу
func ParseParams(params H, form *url.Values) (*url.Values, error) {
	for key, any := range params {
		switch v := any.(type) {
		case string:
			form.Add(key, any.(string))
		case int, int8, int16, int32, int64:
			form.Add(key, fmt.Sprintf("%v", any))
		case float32, float64:
			form.Add(key, fmt.Sprintf("%f", any))
		case []string:
			form.Add(key, strings.Join(any.([]string), ","))
		case []int:
			form.Add(key, tools.SliceOfIntsToString(any.([]int), ","))
		default:
			return nil, fmt.Errorf("не удалось привести параметр %s=%v типа %T к строке", key, any, v)
		}
	}

	return form, nil
}

// Стандартные параметры для всех запросов к VK API
func (b *Bot) getFormData() *url.Values {
	return &url.Values{
		"v":            {b.apiVersion},
		"access_token": {b.token},
	}
}

// Отправка запроса к VK API
//
// Параметры:
//  method: название метода VK API (прим. messages.send)
//  params: параметры запроса
func (b *Bot) RawRequest(method string, params H) ([]byte, error) {
	apiUrl := *b.baseURL
	apiUrl.Path += method
	form, err := ParseParams(params, b.getFormData())

	if err != nil {
		return []byte{}, errors.WithStack(err)
	}
	resp, err := b.client.PostForm(apiUrl.String(), *form)

	if err != nil {
		return []byte{}, errors.Wrapf(err, "ошибка запроса к VK API. Метод %q, параметры %#v", method, params)
	}

	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return []byte{}, errors.Wrap(err, "ошибка чтения тела ответа")
	}

	return raw, err
}

// Отправка запроса к VK API и запись ответа в объект переданной структуры.
// Параметр structure передавать только как указатель!!!
func (b *Bot) requestStruct(method string, params H, structure interface{}) error {
	raw, err := b.RawRequest(method, params)
	if err != nil {
		return errors.WithStack(err)
	}

	err = json.Unmarshal(raw, structure)

	if err != nil {
		return errors.Wrapf(err, "Ошибка при приведении json ответа %q к Go структуре %q",
			raw, reflect.TypeOf(structure),
		)
	}

	return err
}

// Получение LongPoll сервера от VK API
func (b *Bot) GetLongPollServer(wait int) (*LongPollServer, error) {
	resp, err := b.RawRequest("groups.getLongPollServer", H{"group_id": b.groupID})

	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения Long Poll сервера")
	}

	var longPollServerResp struct {
		Response *LongPollServer `json:"response,omitempty"`
		Error    *ApiError       `json:"error,omitempty"`
	}

	err = json.Unmarshal(resp, &longPollServerResp)

	if err != nil {
		fmt.Printf("err %q\n", err)
		return nil, errors.WithStack(err)
	}

	if longPollServerResp.Error != nil {
		return nil, errors.Errorf("VK не отдал Long Poll сервер: %+v", longPollServerResp.Error)
	}

	if longPollServerResp.Response == nil {
		return nil, errors.Errorf("неизвестная ошибка при запросе Long Poll сервера")
	}

	longPollServerResp.Response.Wait = wait

	return longPollServerResp.Response, err
}

// Общий метод отправки сообщений https://vk.com/dev/messages.send
func (b *Bot) SendMessage(opt *SendOptions) (int, error) {
	if err := opt.validate(); err != nil {
		return 0, errors.WithStack(err)
	}

	params, err := opt.toParams()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	var messageResponse struct {
		Response int       `json:"response"`
		Error    *ApiError `json:"error,omitempty"`
	}

	err = b.requestStruct("messages.send", params, &messageResponse)
	if err != nil {
		return 0, errors.Wrap(err, "ошибка при отправке запроса на отправку сообщения")
	}
	if messageResponse.Error != nil {
		return 0, errors.Errorf("ошибка при отправке сообщения: %q", messageResponse.Error.ErrorMsg)
	}

	return messageResponse.Response, nil
}

// Общий метод отправки сообщений нескольким пользователям
func (b *Bot) SendMultipleMessages(opt *SendOptions, userIDs ...int) ([]SendMultipleMessagesResponse, error) {
	if err := opt.validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	params, err := opt.toParams()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var messageResponse struct {
		Response []SendMultipleMessagesResponse `json:"response,omitempty"`
		Error    *ApiError                      `json:"error,omitempty"`
	}

	params["user_ids"] = tools.SliceOfIntsToString(userIDs, ",")
	err = b.requestStruct("messages.send", params, &messageResponse)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка при отправке запроса на отправку сообщения")
	}
	if messageResponse.Error != nil {
		return nil, errors.Errorf("ошибка при отправке сообщения: %q", messageResponse.Error.ErrorMsg)
	}

	return messageResponse.Response, nil
}

// Отправка текстового сообщения
func (b *Bot) SendTextMessage(message string, peerID int) (int, error) {
	opt := &SendOptions{
		PeerID:  peerID,
		Message: message,
	}

	resp, err := b.SendMessage(opt)

	if err != nil {
		return 0, errors.WithStack(err)
	}

	return resp, nil
}

// Отправка текстового сообщения нескольким пользователям
func (b *Bot) SendMultipleTextMessages(message string, userIDs ...int) ([]SendMultipleMessagesResponse, error) {
	opt := &SendOptions{
		Message:        message,
		IsSendMultiple: true,
	}

	resp, err := b.SendMultipleMessages(opt, userIDs...)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return resp, nil
}

// Отправка сообщения с клавиатурой
func (b *Bot) SendKeyboardMessage(message string, keyboard *Keyboard, peerID int) (int, error) {
	opt := &SendOptions{
		PeerID:   peerID,
		Keyboard: keyboard,
		Message:  message,
	}

	resp, err := b.SendMessage(opt)

	if err != nil {
		return 0, errors.WithStack(err)
	}

	return resp, nil
}

// Отправка сообщения с клавиатурой нескольким пользователям
func (b *Bot) SendMultipleKeyboardMessages(
	message string, keyboard *Keyboard, userIDs ...int,
) ([]SendMultipleMessagesResponse, error) {
	opt := &SendOptions{
		Message:        message,
		Keyboard:       keyboard,
		IsSendMultiple: true,
	}

	resp, err := b.SendMultipleMessages(opt, userIDs...)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return resp, nil
}
