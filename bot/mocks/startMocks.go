package mocks

import (
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"time"
)

var StartMessage *vk.MessageNew

func InitStartMocks() {
	StartMessage = &vk.MessageNew{
		Message: &vk.Message{
			ID:                    1176,
			Date:                  1588093139,
			PeerID:                120141017,
			FromID:                120141017,
			Text:                  "Начать",
			ConversationMessageId: 1147,
			Important:             false,
			RandomID:              0,
			IsHidden:              false,
		},
		ClientInfo: vk.ClientInfo{
			ButtonActions:  []string{"text", "vkpay", "open_app", "location", "open_link"},
			Keyboard:       true,
			InlineKeyboard: true,
			Carousel:       false,
			LangID:         0,
		},
		UpdateMeta: &vk.UpdateMeta{
			ReceiveTime: time.Time{},
			UpdateType:  vk.MessageNewUpdate,
			EventID:     "a098902bec84ddb471adf198e871a9a0945fd3c9",
			GroupID:     123,
			Handled:     false,
		},
	}
}

func init() {
	InitStartMocks()
}
