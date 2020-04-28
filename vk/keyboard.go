package vk

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/qulaz/khti_timetable_bot/vk/tools"
)

// Цвета кнопок
const (
	COLOR_PRIMARY   = "primary"
	COLOR_SECONDARY = "secondary"
	COLOR_NEGATIVE  = "negative"
	COLOR_POSITIVE  = "positive"
)

// Типы кнопок
const (
	TEXT_BUTTON     = "text"
	LINK_BUTTON     = "open_link"
	LOCATION_BUTTON = "location"
	VK_PAY_BUTTON   = "vkpay"
	VK_APP_BUTTON   = "open_app"
)

var specialButtons = []string{LOCATION_BUTTON, VK_PAY_BUTTON, VK_APP_BUTTON}

// Общий интерфейс для кнопок
type Button interface {
	GetAction() ButtonAction
}

// Общий интерфейс для ButtonAction
type ButtonAction interface {
	GetType() string
}

type Keyboard struct {
	// https://vk.com/dev/bots_docs_3
	OneTime bool       `json:"one_time"`
	Buttons [][]Button `json:"buttons"`
	Inline  bool       `json:"inline"`

	// Служебные поля
	isLastButtonSpecial bool
}

// Колличество строк в клавиатуре
func (k *Keyboard) RowsCount() int {
	return len(k.Buttons)
}

// Количество кнопок в клавиатуре
func (k *Keyboard) ButtonCount() int {
	count := 0
	for _, row := range k.Buttons {
		count += len(row)
	}

	return count
}

// Дополнительная информация о кнопке. Подробнее об отправке клавиатуры и о Payload: https://vk.cc/a2sn33
type ButtonPayload struct {
	// Команда
	Command string `json:"command"`
	// Параметры команды
	Body   string `json:"body,omitempty"`
	Offset int    `json:"offset,omitempty"`
	// Дополнительные поля
	Other *H `json:"other,omitempty"`
}

type TextButtonAction struct {
	// https://vk.com/dev/bots_docs_3
	Type    string         `json:"type"` // text
	Label   string         `json:"label"`
	Payload *ButtonPayload `json:"payload,omitempty"`
}

func (a TextButtonAction) GetType() string {
	return a.Type
}

type TextButton struct {
	// https://vk.com/dev/bots_docs_3
	Action TextButtonAction `json:"action"`
	// Значения: primary, secondary, negative, positive
	Color string `json:"color"`
}

func (b TextButton) GetAction() ButtonAction {
	return b.Action
}

type LinkButtonAction struct {
	// https://vk.com/dev/bots_docs_3
	Type  string `json:"type"` // open_link
	Link  string `json:"link"`
	Label string `json:"label"`
	// Не используется. Передается для совместимости со старыми клиентами
	Payload string `json:"payload"`
}

func (a LinkButtonAction) GetType() string {
	return a.Type
}

type LinkButton struct {
	// https://vk.com/dev/bots_docs_3
	Action LinkButtonAction `json:"action"`
}

func (b LinkButton) GetAction() ButtonAction {
	return b.Action
}

type LocationButtonAction struct {
	// https://vk.com/dev/bots_docs_3
	Type    string         `json:"type"` // location
	Payload *ButtonPayload `json:"payload,omitempty"`
}

func (a LocationButtonAction) GetType() string {
	return a.Type
}

type LocationButton struct {
	// https://vk.com/dev/bots_docs_3
	Action LocationButtonAction `json:"action"`
}

func (b LocationButton) GetAction() ButtonAction {
	return b.Action
}

type VkPayButtonAction struct {
	// https://vk.com/dev/bots_docs_3
	Type string `json:"type"` // vkpay
	// Не используется. Передается для совместимости со старыми клиентами
	Payload string `json:"payload"`
	Hash    string `json:"hash"`
}

func (a VkPayButtonAction) GetType() string {
	return a.Type
}

type VkPayButton struct {
	// https://vk.com/dev/bots_docs_3
	Action VkPayButtonAction `json:"action"`
}

func (b VkPayButton) GetAction() ButtonAction {
	return b.Action
}

type VkAppsButtonAction struct {
	// https://vk.com/dev/bots_docs_3
	Type    string `json:"type"` // open_app
	AppID   int    `json:"app_id"`
	OwnerID int    `json:"owner_id"`
	// Не используется. Передается для совместимости со старыми клиентами
	Payload string `json:"payload"`
	Label   string `json:"label"`
	Hash    string `json:"hash"`
}

func (a VkAppsButtonAction) GetType() string {
	return a.Type
}

type VkAppsButton struct {
	// https://vk.com/dev/bots_docs_3
	Action VkAppsButtonAction `json:"action"`
}

func (b VkAppsButton) GetAction() ButtonAction {
	return b.Action
}

// Создание объекта клавиатуры
func NewKeyboard(oneTime bool) *Keyboard {
	return &Keyboard{
		OneTime: oneTime,
		Buttons: make([][]Button, 0, 2),
		Inline:  false,
	}
}

// Создание объекта инлайн клавиатуры
func NewInlineKeyboard() *Keyboard {
	return &Keyboard{
		OneTime: false,
		Buttons: make([][]Button, 0, 2),
		Inline:  true,
	}
}

// Добавить ряд кнопок. Функция соблюдает требования VK API по поводу количества рядов
func (k *Keyboard) AddRow() error {
	rowsCount := k.RowsCount()

	if (rowsCount >= 10 && !k.Inline) || (rowsCount >= 6 && k.Inline) {
		return errors.New("Превышено максимальное количество строк")
	}
	if rowsCount > 0 {
		if len(k.Buttons[rowsCount-1]) == 0 {
			return nil
		}
	}
	k.Buttons = append(k.Buttons, make([]Button, 0))
	return nil
}

// Добавлеие кнопки в клавиатуру. Функция соблюдает все требования VK API по поводу размещения и количества кнопок
func (k *Keyboard) AddButton(b Button) error {
	rowsCount := k.RowsCount()
	btnCount := k.ButtonCount()

	// Проверка на превышение кол-ва кнопок
	if (!k.Inline && btnCount == 40) || (k.Inline && btnCount == 10) {
		return errors.New("Достигнуто максимальное количество кнопок в клавиатуре")
	}

	// Добавление первого ряда кнопок
	if rowsCount == 0 {
		if err := k.AddRow(); err != nil {
			return errors.WithStack(err)
		}
		return k.AddButton(b)
	}

	// Предотвращение добавления кнопки в последний ряд, при условии, что предыдущяя кнопка - "специальная"
	if k.isLastButtonSpecial {
		return k.AddRow() // вернет ошибку
	}

	buttonsInRow := len(k.Buttons[rowsCount-1]) // кол-во кнопок в строке

	// "Специальные" кнопки должны занимать всю строку целиком
	if tools.IsStringInSlice(b.GetAction().GetType(), specialButtons) {
		// Если в строке уже есть кнопки - добавляем новую строку для "специальной" кнопки
		if buttonsInRow > 0 {
			if err := k.AddRow(); err != nil {
				return err
			}
			rowsCount = k.RowsCount() // обновление кол-ва рядов
		}
		k.Buttons[rowsCount-1] = append(k.Buttons[rowsCount-1], b)
		// если это не последний из возможных рядов - добавляем после "специальной" кнопки новый ряд
		if (rowsCount < 10 && !k.Inline) || (rowsCount < 6 && k.Inline) {
			return k.AddRow()
		} else {
			// Иначе ставим флаг для предотвращения добавления новых кнопок в один последний ряд со "специальной"
			k.isLastButtonSpecial = true
		}
		return nil
	}

	// Добавление новой строки при переполнении ряда
	if buttonsInRow == 5 {
		if err := k.AddRow(); err != nil {
			return errors.WithStack(err)
		}
	}

	// Добавление кнопки
	k.Buttons[rowsCount-1] = append(k.Buttons[rowsCount-1], b)
	return nil
}

func (k *Keyboard) AddTextButton(label, color string, payload *ButtonPayload) error {
	b := TextButton{
		Action: TextButtonAction{
			Type:    TEXT_BUTTON,
			Label:   label,
			Payload: payload,
		},
		Color: color,
	}
	if err := k.AddButton(b); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (k *Keyboard) AddLinkButton(label, link string) error {
	b := LinkButton{
		Action: LinkButtonAction{
			Type:    LINK_BUTTON,
			Link:    link,
			Label:   label,
			Payload: "",
		},
	}
	if err := k.AddButton(b); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (k *Keyboard) AddLocationButton(payload *ButtonPayload) error {
	b := LocationButton{
		Action: LocationButtonAction{
			Type:    LOCATION_BUTTON,
			Payload: payload,
		},
	}
	if err := k.AddButton(b); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (k *Keyboard) AddVkPayButton(hash string) error {
	b := VkPayButton{
		Action: VkPayButtonAction{
			Type:    VK_PAY_BUTTON,
			Payload: "",
			Hash:    hash,
		},
	}
	if err := k.AddButton(b); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (k *Keyboard) AddVkAppsButton(appID, ownerID int, label, hash string) error {
	b := VkAppsButton{
		Action: VkAppsButtonAction{
			Type:    VK_APP_BUTTON,
			AppID:   appID,
			OwnerID: ownerID,
			Payload: "",
			Label:   label,
			Hash:    hash,
		},
	}
	if err := k.AddButton(b); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Конвертирование клавиатуры в json-строку
func (k *Keyboard) ToJson() (string, error) {
	b, err := json.Marshal(k)

	if err != nil {
		return "", errors.Wrap(err, "ошибка конвертирования клавиатуры в json")
	}

	return string(b), err
}
