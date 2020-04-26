package vk

import (
	"gitlab.com/qulaz/khti_timetable_bot/vk/tools"
)

// Информация о пользователе, приходящяя в обновлении message_new
type ClientInfo struct {
	ButtonActions  []string `json:"button_actions"`
	Keyboard       bool     `json:"keyboard"`
	InlineKeyboard bool     `json:"inline_keyboard"`
	Carousel       bool     `json:"carousel"`
	LangID         int      `json:"lang_id"`
}

// Есть ли поддержка клавиатуры у клиента
func (c *ClientInfo) IsKeyboardSupported() bool {
	return c.Keyboard
}

// Поддерживаются ли текстовые кнопки у клиента
func (c *ClientInfo) IsTextButtonSupported() bool {
	return tools.IsStringInSlice("text", c.ButtonActions)
}

// Поддерживаются ли кнопки с ссылками у клиента
func (c *ClientInfo) IsLinkButtonSupported() bool {
	return tools.IsStringInSlice("open_link", c.ButtonActions)
}

// Поддерживаются ли кнопки VK Pay у клиента
func (c *ClientInfo) IsVkPayButtonSupported() bool {
	return tools.IsStringInSlice("vkpay", c.ButtonActions)
}

// Поддерживаются ли кнопки локации у клиента
func (c *ClientInfo) IsLocationButtonSupported() bool {
	return tools.IsStringInSlice("location", c.ButtonActions)
}

// Поддерживаются ли VK App кнопки у клиента
func (c *ClientInfo) IsVkAppButtonSupported() bool {
	return tools.IsStringInSlice("open_app", c.ButtonActions)
}

// Есть ли поддержка инлайн клавиатуры у клиента
func (c *ClientInfo) IsInlineKeyboardSupported() bool {
	return c.InlineKeyboard
}

// Есть ли поддержка карусели у клиента
func (c *ClientInfo) IsCarouselSupported() bool {
	return c.Carousel
}
