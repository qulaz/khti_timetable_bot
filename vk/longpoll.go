package vk

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/qulaz/khti_timetable_bot/vk/tools"
	"net/url"
	"strconv"
)

const (
	DEFAULT_LONG_POLLING_WAIT = 25
	DEFAULT_REQUEST_INTERVAL  = 400
)

type LongPollServer struct {
	// https://vk.com/dev/groups.getLongPollServer
	Key    string `json:"key"`
	Server string `json:"server"`
	Ts     string `json:"ts"`

	Wait int
}

type LongPollUpdate struct {
	Type    string `json:"type"`
	Object  H      `json:"object"`
	GroupID int    `json:"group_id"`
	EventID string `json:"event_id"`
}

type LongPollResponse struct {
	Ts      json.Number      `json:"ts"`
	Updates []LongPollUpdate `json:"updates"`
	Failed  int              `json:"failed"`
}

// Собираем URL-адрес для запроса из объекта Long Poll сервера
func (s *LongPollServer) buildUrl() (*url.URL, error) {
	base, err := url.Parse(s.Server)

	if err != nil {
		return nil, errors.Wrapf(err, "Ошибка парсинга URL Long Poll сервера %#v", s.Server)
	}

	// Query params
	params := url.Values{}
	params.Add("act", "a_check")
	params.Add("key", s.Key)
	params.Add("ts", s.Ts)
	params.Add("wait", strconv.Itoa(s.Wait))
	base.RawQuery = params.Encode()

	return base, err
}

// Получение и подготовка к обработке апдейтов от Long Poll сервера
func GetUpdates(b *Bot, server *LongPollServer) ([]*Update, error) {
	serverUrl, err := server.buildUrl()
	if err != nil {
		return []*Update{}, errors.WithStack(err)
	}

	resp, err := getRequest(b.client, serverUrl.String())
	if err != nil {
		return []*Update{}, errors.WithStack(err)
	}

	response := &LongPollResponse{}
	if err := json.Unmarshal(resp, response); err != nil {
		return []*Update{}, errors.Wrap(err, "bad json")
	}

	// Если вернулась ошибка - обновляем сервер
	if response.Failed > 0 {
		newServer, err := b.GetLongPollServer(server.Wait)
		if err != nil {
			return nil, errors.Wrap(err, "ошибка при получении нового Long Poll сервера")
		}
		server.Ts = newServer.Ts
		server.Server = newServer.Server
		server.Key = newServer.Key
		server.Wait = newServer.Wait

		return []*Update{}, nil
	}

	if response.Ts != "" {
		server.Ts = string(response.Ts)
	} else {
		b.Logger.Warnf("В ответе Long Poll сервера нет следующего ts: %+v; raw=%+v", response, resp)
	}

	if l := len(response.Updates); l > 0 {
		updates := make([]*Update, 0, l)

		for _, upd := range response.Updates {
			updateMap := make(H, 1)
			updateObj := &Update{
				UpdateMeta: &UpdateMeta{
					ReceiveTime: tools.Now(),
					UpdateType:  upd.Type,
					EventID:     upd.EventID,
					GroupID:     upd.GroupID,
				},
			}

			// структуа поля "object" апдейтов "message_edit" и "message_reply" совпадает со структурой объекта сообщения
			// поэтому, чтобы избижать дублировнания, оборачивам его содержимое в дополниельное поле с ключем "object",
			// которое затем, в структурах этих апдейтов парсим
			if upd.Type == MessageEditUpdate || upd.Type == MessageReplyUpdate {
				upd.Object = map[string]interface{}{"object": upd.Object}
			}

			updateMap[upd.Type] = upd.Object
			if err := tools.MapToStruct(updateMap, updateObj); err != nil {
				b.Logger.Errorf("Ошибка десериализации объекта ивента %q %+v: %q", upd.Type, upd.Object, err)
				continue
			}

			updates = append(updates, updateObj)
		}

		return updates, nil
	}

	return []*Update{}, nil
}
