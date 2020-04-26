package vk

import "time"

// Константы типов обновлений
const (
	MessageNewUpdate   = "message_new"
	MessageEditUpdate  = "message_edit"
	MessageReplyUpdate = "message_reply"
	MessageAllowUpdate = "message_allow"
	MessageDenyUpdate  = "message_deny"
)

type Update struct {
	MessageNew   *MessageNew   `json:"message_new"`
	MessageEdit  *MessageEdit  `json:"message_edit"`
	MessageReply *MessageReply `json:"message_reply"`
	MessageAllow *MessageAllow `json:"message_allow"`
	MessageDeny  *MessageDeny  `json:"message_deny"`

	// Служебная информация об апдейте
	UpdateMeta *UpdateMeta `json:"-"`
}

type UpdateMeta struct {
	// Время получения апдейта в UTC
	ReceiveTime time.Time `json:"-"`
	// Тип апдейта
	UpdateType string `json:"-"`
	// Идентификатор апдейта
	EventID string `json:"-"`
	// Групп, к которой относится этот апдейт
	GroupID int `json:"-"`
	// Статус обработки апдейта
	Handled bool `json:"-"`
}

// Интерфейс для конкретного типа обновлений (MessageNew, MessageEdit и тд.)
type Updater interface {
	GetType() string
	SetUpdateMeta(meta *UpdateMeta)
}
