package mocks

import (
	"fmt"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"gitlab.com/qulaz/khti_timetable_bot/vk/tests/consts"
	"time"
)

var (
	HandleMessageText    = "Начало"
	HandleMessageNewText = "Обычный текст сообщения"
	HandleCommandText    = "/команда тело команды"
	HandleCommandCommand = "/команда"
	HandleCommandBody    = "тело команды"
)
var TestTime = time.Date(2020, 3, 25, 15, 15, 15, 0, time.UTC)
var (
	UpdateMessageReply      *vk.Update
	UpdateMessageNew        *vk.Update
	UpdateMessageEdit       *vk.Update
	UpdateMessageAllow      *vk.Update
	UpdateMessageDeny       *vk.Update
	Update_HandleMessage    *vk.Update
	Update_HandleMessageNew *vk.Update
	Update_HandleCommand    *vk.Update
	Update_HandleCommand2   *vk.Update
)
var LongPollTestUpdates [][]interface{}

func ResetMocks() {
	UpdateMessageReply = &vk.Update{
		MessageReply: &vk.MessageReply{
			Message: &vk.Message{
				ID:                    392,
				Date:                  1583332586,
				PeerID:                123,
				FromID:                -123,
				Text:                  "l",
				RandomID:              8494866,
				Important:             false,
				ConversationMessageId: 363,
				IsHidden:              false,
				Out:                   1,
			},
		},
		UpdateMeta: &vk.UpdateMeta{
			ReceiveTime: TestTime,
			UpdateType:  vk.MessageReplyUpdate,
			EventID:     "03428f0338eb60f6e397c10c4739f8274b770a18",
			GroupID:     123,
		},
	}
	UpdateMessageNew = &vk.Update{
		MessageNew: &vk.MessageNew{
			Message: &vk.Message{
				ID:                    387,
				Date:                  1583246250,
				PeerID:                123,
				FromID:                -123,
				Text:                  "kek",
				RandomID:              0,
				Important:             false,
				ConversationMessageId: 358,
				IsHidden:              false,
				Out:                   0,
			},
			ClientInfo: vk.ClientInfo{
				ButtonActions:  []string{"text", "vkpay", "open_app", "location", "open_link"},
				Keyboard:       true,
				InlineKeyboard: true,
				Carousel:       false,
				LangID:         0,
			},
		},
		UpdateMeta: &vk.UpdateMeta{
			ReceiveTime: TestTime,
			UpdateType:  vk.MessageNewUpdate,
			EventID:     "fe989fe0db287324f53df8c9454876ee7a6752f1",
			GroupID:     123,
		},
	}
	UpdateMessageEdit = &vk.Update{
		MessageEdit: &vk.MessageEdit{
			Message: &vk.Message{
				Date:                  1583338337,
				FromID:                123,
				ID:                    364,
				Out:                   1,
				PeerID:                123,
				Text:                  "лол",
				ConversationMessageId: 364,
				RandomID:              0,
				UpdateTime:            1583338337,
			},
		},
		UpdateMeta: &vk.UpdateMeta{
			ReceiveTime: TestTime,
			UpdateType:  vk.MessageEditUpdate,
			EventID:     "16aa8578251c1b336701b8defc7f98c22f62d804",
			GroupID:     123,
		},
	}
	UpdateMessageAllow = &vk.Update{
		MessageAllow: &vk.MessageAllow{
			UserID: 123,
			Key:    "testKey",
		},
		UpdateMeta: &vk.UpdateMeta{
			ReceiveTime: TestTime,
			UpdateType:  vk.MessageAllowUpdate,
			EventID:     "16aa8578251c1b336701b8defc7f98c22f62d804",
			GroupID:     123,
		},
	}
	UpdateMessageDeny = &vk.Update{
		MessageDeny: &vk.MessageDeny{UserID: 123},
		UpdateMeta: &vk.UpdateMeta{
			ReceiveTime: TestTime,
			UpdateType:  vk.MessageDenyUpdate,
			EventID:     "16aa8578251c1b336701b8defc7f98c22f62d804",
			GroupID:     123,
		},
	}
	Update_HandleMessage = &vk.Update{
		MessageNew: &vk.MessageNew{
			Message: &vk.Message{
				ID:     1,
				Date:   1515151515,
				PeerID: 1,
				FromID: 1,
				Text:   HandleMessageText,
			},
			ClientInfo: vk.ClientInfo{},
		},
		UpdateMeta: &vk.UpdateMeta{
			ReceiveTime: time.Now(), UpdateType: vk.MessageNewUpdate, EventID: "asfas5f1as56f",
			GroupID: consts.TEST_GROUP_ID,
		},
	}
	Update_HandleMessageNew = &vk.Update{
		MessageNew: &vk.MessageNew{
			Message: &vk.Message{
				ID:     1,
				Date:   1515151515,
				PeerID: 1,
				FromID: 1,
				Text:   HandleMessageNewText,
			},
			ClientInfo: vk.ClientInfo{},
		},
		UpdateMeta: &vk.UpdateMeta{
			ReceiveTime: time.Now(), UpdateType: vk.MessageNewUpdate, EventID: "asfas5f1as56f",
			GroupID: consts.TEST_GROUP_ID,
		},
	}
	Update_HandleCommand = &vk.Update{
		MessageNew: &vk.MessageNew{
			Message: &vk.Message{
				ID:     1,
				Date:   1515151515,
				PeerID: 1,
				FromID: 1,
				Text:   HandleCommandText,
			},
			ClientInfo: vk.ClientInfo{},
		},
		UpdateMeta: &vk.UpdateMeta{
			ReceiveTime: time.Now(), UpdateType: vk.MessageNewUpdate, EventID: "asfas5f1as56f",
			GroupID: consts.TEST_GROUP_ID,
		},
	}
	Update_HandleCommand2 = &vk.Update{
		MessageNew: &vk.MessageNew{
			Message: &vk.Message{
				ID:     1,
				Date:   1515151515,
				PeerID: 1,
				FromID: 1,
				Text:   "Команда",
				RawPayload: fmt.Sprintf(
					"{\"command\": \"%s\", \"body\": \"%s\"}", HandleCommandCommand, HandleCommandBody,
				),
			},
			ClientInfo: vk.ClientInfo{},
		},
		UpdateMeta: &vk.UpdateMeta{
			ReceiveTime: time.Now(), UpdateType: vk.MessageNewUpdate, EventID: "asfas5f1as56f",
			GroupID: consts.TEST_GROUP_ID,
		},
	}

	LongPollTestUpdates = [][]interface{}{
		{UpdateMessageReply, UpdateMessageNew},
		{UpdateMessageEdit, UpdateMessageAllow, UpdateMessageDeny},
	}
}

func init() {
	ResetMocks()
}
