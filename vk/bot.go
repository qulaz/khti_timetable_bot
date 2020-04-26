package vk

import (
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"reflect"
)

const (
	handleMessageCode = "\amsg"
	handleCommandCode = "\acommand"
)

type Bot struct {
	groupID    int
	token      string
	apiVersion string
	baseURL    *url.URL

	Logger Logger
	client *http.Client
	poller Poller

	handlers    map[string]interface{}
	middlewares map[string][]func(b *Bot, u interface{}) bool
}

// Опции для настройки бота
type Settings struct {
	// Обязательные параметры
	//
	// Идентификатор группы в которой будет работать бот
	GroupID int
	// Токен бота
	Token string

	// Необязательные параметры
	//
	// По умолчанию используется LongPoller
	Poller Poller
	// Объект логгера, реализующий одноименный интерфейс
	Logger Logger
	// Настроенный HTTP Client для отправки запросов к VK API
	Client *http.Client
	// Базовый URL VK API
	BaseURL string
	// Версия VK API
	ApiVersion string
}

// Функция создания объекта бота
func CreateBot(s Settings) (*Bot, error) {
	if s.Token == "" {
		return nil, errors.New("Token является обязательным параметром")
	}
	if s.GroupID == 0 {
		return nil, errors.New("GroupID является обязательным параметром")
	}
	if s.BaseURL == "" {
		s.BaseURL = VK_API_BASE_URL
	}
	if s.ApiVersion == "" {
		s.ApiVersion = VK_API_VERISON
	}
	baseUrl, err := url.Parse(s.BaseURL)
	if err != nil {
		return nil, errors.Errorf("Ошибка при парсинге BaseURL: %v", err)
	}

	poller := s.Poller
	if poller == nil {
		poller = NewDefaultLongPoller()
	}

	client := s.Client
	if client == nil {
		client = NewClient()
	}

	logger := s.Logger
	if logger == nil {
		logger = &DefaultLogType{}
	}

	bot := &Bot{
		groupID:     s.GroupID,
		token:       s.Token,
		handlers:    make(map[string]interface{}),
		client:      client,
		Logger:      logger,
		apiVersion:  s.ApiVersion,
		poller:      poller,
		baseURL:     baseUrl,
		middlewares: make(map[string][]func(b *Bot, u interface{}) bool),
	}

	return bot, nil
}

func (b *Bot) SetBaseUrl(baseUrl string) error {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return errors.Wrap(err, "ошибка при установке BaseURL")
	}
	b.baseURL = u

	return nil
}

// Добавление обработчика, который будет срабатывать в случае, если придет сообщение полностью идентичное заданному
func (b *Bot) HandleMessage(message string, handler HandlerMessageNew) {
	b.handlers[handleMessageCode+message] = handler
}

// Добавление обработчика, который срабатывает на пришедшие комманды. Командами считаются сообщения, которые начинаются
// с символа /, а также содержимое поля Command в ButtonPayload кнопок.
func (b *Bot) HandleCommand(command string, handler HandlerMessageNew) {
	b.handlers[handleCommandCode+command] = handler
}

// Добавление нового обработчика, который срабатывает на все пришедшие сообщения, не подошедшие под критерии остальных
// обработчиков
func (b *Bot) HandleMessageNew(handler HandlerMessageNew) {
	b.handlers[MessageNewUpdate] = handler
}

// Добавление обработчика на событие message_edit
func (b *Bot) HandleMessageEdit(handler HandlerMessageEdit) {
	b.handlers[MessageEditUpdate] = handler
}

// Добавление обработчика на событие message_reply
func (b *Bot) HandleMessageReply(handler HandlerMessageReply) {
	b.handlers[MessageReplyUpdate] = handler
}

// Добавление обработчика на событие message_deny
func (b *Bot) HandleMessageDeny(handler HandlerMessageDeny) {
	b.handlers[MessageDenyUpdate] = handler
}

// Добавление обработчика на событие message_allow
func (b *Bot) HandleMessageAllow(handler HandlerMessageAllow) {
	b.handlers[MessageAllowUpdate] = handler
}

// Добавление Middleware
func (b *Bot) AddMiddleware(middleware Middleware) {
	v := reflect.ValueOf(middleware)

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsNil() {
			middlewareKey := v.Type().Field(i).Name

			if _, ok := b.middlewares[v.Type().Field(i).Name]; !ok {
				b.middlewares[middlewareKey] = make([]func(b *Bot, u interface{}) bool, 0, 2)
			}

			switch n := middlewareKey; n {
			case onPreProcessMiddleware, onPostProcessMiddleware:
				fn := v.Field(i).Interface().(func(b *Bot, u *Update) bool)
				b.middlewares[n] = append(b.middlewares[n], func(b *Bot, u interface{}) bool {
					return fn(b, u.(*Update))
				})
			case onPreMessageNewMiddleware, onPostMessageNewMiddleware:
				fn := v.Field(i).Interface().(func(b *Bot, u *MessageNew) bool)
				b.middlewares[n] = append(b.middlewares[n], func(b *Bot, u interface{}) bool {
					return fn(b, u.(*MessageNew))
				})
			case onPreMessageReplyMiddleware, onPostMessageReplyMiddleware:
				fn := v.Field(i).Interface().(func(b *Bot, u *MessageReply) bool)
				b.middlewares[n] = append(b.middlewares[n], func(b *Bot, u interface{}) bool {
					return fn(b, u.(*MessageReply))
				})
			case onPreMessageEditMiddleware, onPostMessageEditMiddleware:
				fn := v.Field(i).Interface().(func(b *Bot, u *MessageEdit) bool)
				b.middlewares[n] = append(b.middlewares[n], func(b *Bot, u interface{}) bool {
					return fn(b, u.(*MessageEdit))
				})
			case onPreMessageAllowMiddleware, onPostMessageAllowMiddleware:
				fn := v.Field(i).Interface().(func(b *Bot, u *MessageAllow) bool)
				b.middlewares[n] = append(b.middlewares[n], func(b *Bot, u interface{}) bool {
					return fn(b, u.(*MessageAllow))
				})
			case onPreMessageDenyMiddleware, onPostMessageDenyMiddleware:
				fn := v.Field(i).Interface().(func(b *Bot, u *MessageDeny) bool)
				b.middlewares[n] = append(b.middlewares[n], func(b *Bot, u interface{}) bool {
					return fn(b, u.(*MessageDeny))
				})
			}
		}
	}
}

// Запуск бота
func (b *Bot) Run() error {
	err := b.poller.Poll(b)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Функция обработки обновлений
func (b *Bot) ProcessUpdates(updates []*Update) {
	for _, update := range updates {
		if u := update.MessageNew; u != nil {
			u.Message.ProcessMessage()

			handlerKeys := []string{
				handleCommandCode + u.Message.MessageCommand, handleMessageCode + u.Message.Text, MessageNewUpdate,
			}

			if b.handle(u, update, onPreMessageNewMiddleware, onPostMessageNewMiddleware, handlerKeys...) {
				continue
			}
		}
		if u := update.MessageReply; u != nil {
			u.Message.ProcessMessage()

			if b.handle(u, update, onPreMessageReplyMiddleware, onPostMessageReplyMiddleware, MessageReplyUpdate) {
				continue
			}
		}
		if u := update.MessageEdit; u != nil {
			u.Message.ProcessMessage()

			if b.handle(u, update, onPreMessageEditMiddleware, onPostMessageEditMiddleware, MessageEditUpdate) {
				continue
			}
		}
		if u := update.MessageAllow; u != nil {
			if b.handle(u, update, onPreMessageAllowMiddleware, onPostMessageAllowMiddleware, MessageAllowUpdate) {
				continue
			}
		}
		if u := update.MessageDeny; u != nil {
			if b.handle(u, update, onPreMessageDenyMiddleware, onPostMessageDenyMiddleware, MessageDenyUpdate) {
				continue
			}
		}
	}
}

// Шорткат для запуска всех мидлварей определенного типа
func (b *Bot) triggerMiddleware(u interface{}, middlewareKey string) bool {
	if middlewares, ok := b.middlewares[middlewareKey]; ok {
		for _, middleware := range middlewares {
			if !middleware(b, u) {
				return false
			}
		}
	}

	return true
}

// Обработка обновления
func (b *Bot) handle(u Updater, update *Update, preMiddleware, postMiddleware string, handlerKeys ...string) bool {
	u.SetUpdateMeta(update.UpdateMeta)

	if !b.triggerMiddleware(update, onPreProcessMiddleware) || !b.triggerMiddleware(u, preMiddleware) {
		return false
	}

	for _, handlerKey := range handlerKeys {
		if b.handleUpdate(handlerKey, u) {
			update.UpdateMeta.Handled = true
			break
		}
	}

	b.triggerMiddleware(u, postMiddleware)
	b.triggerMiddleware(update, onPostProcessMiddleware)

	return true
}

// Обработка обновления
func (b *Bot) handleUpdate(key string, update Updater) bool {
	if handler, ok := b.handlers[key]; ok {
		switch handler.(type) {
		case HandlerMessageNew:
			if u, ok := update.(*MessageNew); ok {
				handler.(HandlerMessageNew)(b, u)
				return true
			}
		case HandlerMessageEdit:
			if u, ok := update.(*MessageEdit); ok {
				handler.(HandlerMessageEdit)(b, u)
				return true
			}
		case HandlerMessageReply:
			if u, ok := update.(*MessageReply); ok {
				handler.(HandlerMessageReply)(b, u)
				return true
			}
		case HandlerMessageAllow:
			if u, ok := update.(*MessageAllow); ok {
				handler.(HandlerMessageAllow)(b, u)
				return true
			}
		case HandlerMessageDeny:
			if u, ok := update.(*MessageDeny); ok {
				handler.(HandlerMessageDeny)(b, u)
				return true
			}
		}
	}

	return false
}
