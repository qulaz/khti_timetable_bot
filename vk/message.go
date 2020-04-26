package vk

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/qulaz/khti_timetable_bot/vk/tools"
	"math/rand"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// Regexp для поиска упоминаний бота в тексте сообщения: [club11111|@short_link]
const RE_REMOVE_BOT_IDENTIFIER = `\[club\d+\|[^\]\)]+]`

// Regexp для поиска команд в тексте сообщения: </command> <body>
const RE_MESSAGE_WITH_COMMAND = `(?P<command>/\S+)(?P<whitespace> )?(?P<body>.+)?`

var (
	// Скомпилированный Regexp для поиска и удаления упоминаний бота из текста сообщения
	ReRemoveBotIdentifier = regexp.MustCompile(RE_REMOVE_BOT_IDENTIFIER)
	// Скомпилированный Regexp для поиска команд в тексте сообщения
	ReMessageWithCommand = regexp.MustCompile(RE_MESSAGE_WITH_COMMAND)
)

type Message struct {
	// https://vk.com/dev/objects/message
	//
	ID            int            `json:"id"`
	Date          int            `json:"date"`
	PeerID        int            `json:"peer_id"`
	FromID        int            `json:"from_id"`
	Text          string         `json:"text"`
	RandomID      int64          `json:"random_id"`
	Ref           string         `json:"ref"`
	RefSource     string         `json:"ref_source"`
	Attachments   []interface{}  `json:"attachments"`
	Important     bool           `json:"important"`
	Geo           *Geo           `json:"geo,omitempty"`
	Payload       *ButtonPayload `json:"-"`
	Keyboard      *Keyboard      `json:"keyboard,omitempty"`
	FwdMessages   []*Message     `json:"fwd_messages,omitempty"`
	ReplayMessage *Message       `json:"reply_message,omitempty"`
	Action        *Action        `json:"action,omitempty"`

	// Поля не указанные в документации, но присутствующие в ответе Long Poll сервера
	//
	ConversationMessageId int  `json:"conversation_message_id"`
	IsHidden              bool `json:"is_hidden"`
	Out                   int  `json:"out"`
	UpdateTime            int  `json:"update_time"`     // Только в апдейте message_edit
	AdminAuthorId         int  `json:"admin_author_id"` // В апдейтах message_reply, message_edit

	// Поля обрабатываемые библиотекой
	//
	// Название команды. Берется либо из Payload.Command, либо из текста сообщения.
	// Например из сообщения: "/add такси 250" в MessageCommand попадет "/add"
	MessageCommand string `json:"-"`

	// Аргументы команды. Берется либо из Payload.Body, либо из текста сообщения.
	// Например из сообщения: "/add такси 250" в MessageBody попадет "такси 250"
	MessageBody string `json:"-"`

	// JSON-строка с пейлоадом сообщения
	RawPayload string `json:"payload,omitempty"`
}

// Убираем упоминание бота из текста сообщения и обрезаем возможные лишние пробелы после него
func (m *Message) removeBotIdentifier() {
	m.Text = strings.Trim(ReRemoveBotIdentifier.ReplaceAllString(m.Text, ""), " ")
}

// Создание шорткатов MessageCommand и MessageBody, значения которых берутся из Payload, либо из текста сообщения
func (m *Message) parsePayload() {
	// Преобразование RawPayload (должен быть JSON-строкой!) к ButtonPayload структуре.
	if m.RawPayload != "" {
		p := &ButtonPayload{}
		if err := json.Unmarshal([]byte(m.RawPayload), p); err == nil {
			m.Payload = p
		}
	}

	// Создание шорткатов для названия команды и аргументов команды из Payload
	if m.Payload != nil {
		m.MessageBody = m.Payload.Body
		m.MessageCommand = m.Payload.Command
	}

	// Проверка на наличие команды в тексте сообщения
	match := ReMessageWithCommand.FindAllStringSubmatch(m.Text, -1)
	if len(match) > 0 {
		// "</command> <body>"
		m.MessageCommand, m.MessageBody = match[0][1], match[0][3]
	}
}

// Возвращает true в случае, если сообщение пришло из беседы.
func (m *Message) IsChat() bool {
	if id := m.ID - 2000000000; id > 0 {
		return true
	}

	return false
}

// Обработка сообщения, котрая вклчает в себя:
//  - Удаление из текста сообщения возможного обращения к боту через @/* (которое в апдейте превращается в специальную
//  разметку)
//  - Приведение поля "payload" из апдейта, которое должно являться json-строкой, к ButtonPayload структуре
//  - Создание шорткатов MessageCommand и MessageBody для названия команды и возможных аргументов команды. Данные
//  берутся из пришедшего Payload, либо из текста сообщения
func (m *Message) ProcessMessage() {
	m.removeBotIdentifier()
	m.parsePayload()
}

type Geo struct {
	// https://vk.com/dev/objects/message
	Type        string           `json:"type"`
	Coordinates []GeoCoordinates `json:"coordinates"`
	Place       GeoPlace         `json:"place"`
}

type GeoPlace struct {
	// https://vk.com/dev/objects/message
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Created   int     `json:"created"`
	Icon      string  `json:"icon"`
	Country   string  `json:"country"`
	City      string  `json:"city"`
}

type GeoCoordinates struct {
	// https://vk.com/dev/objects/message
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Action struct {
	// https://vk.com/dev/objects/message
	Type     string `json:"type"`
	MemberID int    `json:"member_id"`
	Text     string `json:"text"`
	Email    string `json:"email"`
	Photo    struct {
		Photo50  string `json:"photo_50"`
		Photo100 string `json:"photo_100"`
		Photo200 string `json:"photo_200"`
	} `json:"photo"`
}

// Параметры для отправки сообщения
type SendOptions struct {
	// https://vk.com/dev/messages.send
	UserID          int       `json:"user_id,omitempty"`
	PeerID          int       `json:"peer_id,omitempty"`
	Domain          string    `json:"domain,omitempty"`
	ChatID          int       `json:"chat_id,omitempty"`
	RandomID        int64     `json:"random_id,omitempty"`
	Message         string    `json:"message,omitempty"`
	Lat             *float64  `json:"lat,omitempty"`
	Long            *float64  `json:"long,omitempty"`
	Attachment      []string  `json:"-"`
	ReplyTo         int       `json:"reply_to,omitempty"`
	ForwardMessages []int     `json:"-"`
	StickerID       int       `json:"sticker_id,omitempty"`
	Keyboard        *Keyboard `json:"-"`
	Payload         *H        `json:"payload,omitempty"`
	DontParseLinks  int       `json:"dont_parse_links,omitempty"`
	DisableMentions int       `json:"disable_mentions,omitempty"`
	Intent          string    `json:"intent,omitempty"`

	// "Сырые" строки с данными, готовые для отправки на сервера VK
	ForwardMessagesProcessed string `json:"forward_messages,omitempty"`
	AttachmentProcessed      string `json:"attachment,omitempty"`
	KeyboardProcessed        string `json:"keyboard,omitempty"`

	// Флаг, указывающий на то, что данное сообщение будет отправляться нескольким получателям сразу
	IsSendMultiple bool `json:"-"`
}

func (opt *SendOptions) validate() error {
	if opt.PeerID == 0 && opt.UserID == 0 && opt.ChatID == 0 && opt.Domain == "" && !opt.IsSendMultiple {
		return errors.New("не задан получатель сообщения")
	}
	if opt.Message == "" && (opt.Lat == nil && opt.Long == nil) && opt.Attachment == nil && opt.ForwardMessages == nil && opt.StickerID == 0 {
		return errors.New("не задан контент сообщения")
	}
	if (opt.Lat != nil && opt.Long == nil) || (opt.Lat == nil && opt.Long != nil) {
		return errors.New("должна быть задана как широта (lat), так и долгота (long)")
	}

	if opt.ForwardMessages != nil {
		opt.ForwardMessagesProcessed = tools.SliceOfIntsToString(opt.ForwardMessages, ",")
	}
	if opt.Attachment != nil {
		opt.AttachmentProcessed = strings.Join(opt.Attachment, ",")
	}
	if opt.Keyboard != nil {
		k, err := opt.Keyboard.ToJson()
		if err != nil {
			return errors.WithStack(err)
		}
		opt.KeyboardProcessed = k
	}
	if opt.RandomID == 0 {
		opt.RandomID = opt.randomID()
	}

	return nil
}

func (opt *SendOptions) getReceiver() int {
	if opt.PeerID != 0 {
		return opt.PeerID
	}
	if opt.UserID != 0 {
		return opt.UserID
	}
	if opt.ChatID != 0 {
		return opt.ChatID
	}
	if opt.Domain != "" {
		return rand.Intn(9999999)
	}
	return int(rand.Int31())
}

func (opt *SendOptions) randomID() int64 {
	return time.Now().UnixNano() + int64(opt.getReceiver()) - int64(rand.Int31())
}

func (opt *SendOptions) toParams() (H, error) {
	params := H{}

	val := reflect.ValueOf(*opt)
	for i := 0; i < val.Type().NumField(); i++ {
		t := val.Type().Field(i)
		field := val.Field(i)
		fieldName := t.Name

		if jsonTag := t.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {
				fieldName = jsonTag[:commaIdx]
				fieldType := field.Type()
				var fieldValue string

				if fieldType.Kind() == reflect.Ptr {
					if elem := field.Elem(); elem.IsValid() {
						switch elem.Interface().(type) {
						case float32, float64:
							fieldValue = fmt.Sprintf("%f", elem.Float())
						}
						params[fieldName] = fieldValue
					}
				} else {
					switch i := field.Interface(); i.(type) {
					case string:
						fieldValue = i.(string)
					case int, int64:
						if v := fmt.Sprintf("%v", field.Int()); v != "0" {
							fieldValue = v
						}
					}

					if fieldValue != "" {
						params[fieldName] = fieldValue
					}

				}
			}
		}
	}

	return params, nil
}

func CreateFloat64(f float64) *float64 {
	return &f
}

type SendMultipleMessagesResponse struct {
	PeerID    int `json:"peer_id"`
	MessageID int `json:"message_id"`
	Error     *struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
	} `json:"error,omitempty"`
}
