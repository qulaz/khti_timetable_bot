package tests

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"gitlab.com/qulaz/khti_timetable_bot/vk/tests/consts"
	"gitlab.com/qulaz/khti_timetable_bot/vk/tests/mocks"
	"os"
	"testing"
)

func init() {
	err := os.Chdir(consts.GetTestDir())
	if err != nil {
		panic(err)
	}
}

func TestBot_ProcessUpdates(t *testing.T) {
	b, _ := vk.CreateBot(vk.Settings{
		GroupID: 1,
		Token:   "TEST_TOKEN",
	})

	var (
		messageNewMessageHandled = false
		messageNewCommandHandled = false
		messageNewHandled        = false
		messageEditHandled       = false
		messageReplyHandled      = false
		messageAllowHandled      = false
		messageDenyHandled       = false
	)

	updates := []*vk.Update{
		mocks.Update_HandleCommand, mocks.Update_HandleMessage, mocks.Update_HandleMessageNew, mocks.UpdateMessageEdit,
		mocks.UpdateMessageReply, mocks.UpdateMessageAllow, mocks.UpdateMessageDeny, mocks.Update_HandleCommand2,
	}

	b.HandleMessageNew(func(b *vk.Bot, m *vk.MessageNew) {
		messageNewHandled = true
		assert.Equal(t, mocks.HandleMessageNewText, m.Message.Text)
	})
	b.HandleMessage(mocks.HandleMessageText, func(b *vk.Bot, m *vk.MessageNew) {
		messageNewMessageHandled = true
		assert.Equal(t, mocks.HandleMessageText, m.Message.Text)
	})
	b.HandleCommand(mocks.HandleCommandCommand, func(b *vk.Bot, m *vk.MessageNew) {
		messageNewCommandHandled = true
		assert.Equal(t, mocks.HandleCommandBody, m.Message.MessageBody)
		assert.Equal(t, mocks.HandleCommandCommand, m.Message.MessageCommand)
	})
	b.HandleMessageEdit(func(b *vk.Bot, m *vk.MessageEdit) {
		messageEditHandled = true
	})
	b.HandleMessageReply(func(b *vk.Bot, m *vk.MessageReply) {
		messageReplyHandled = true
	})
	b.HandleMessageAllow(func(b *vk.Bot, m *vk.MessageAllow) {
		messageAllowHandled = true
	})
	b.HandleMessageDeny(func(b *vk.Bot, m *vk.MessageDeny) {
		messageDenyHandled = true
	})

	b.ProcessUpdates(updates)

	assert.True(t, messageNewMessageHandled)
	assert.True(t, messageNewCommandHandled)
	assert.True(t, messageNewHandled)
	assert.True(t, messageEditHandled)
	assert.True(t, messageReplyHandled)
	assert.True(t, messageAllowHandled)
	assert.True(t, messageDenyHandled)

	for i, u := range updates {
		assert.True(t, u.UpdateMeta.Handled, i)
	}
}

func TestBot_ProcessUpdates_Middleware(t *testing.T) {
	var (
		update_rejected          = false
		OnPreProcessUpdate       = false
		OnPostProcessUpdate      = false
		OnPreMessageNewUpdate    = false
		OnPostMessageNewUpdate   = false
		OnPreMessageReplyUpdate  = false
		OnPostMessageReplyUpdate = false
		OnPreMessageEditUpdate   = false
		OnPostMessageEditUpdate  = false
		OnPreMessageAllowUpdate  = false
		OnPostMessageAllowUpdate = false
		OnPreMessageDenyUpdate   = false
		OnPostMessageDenyUpdate  = false
	)

	b, _ := vk.CreateBot(vk.Settings{
		GroupID: 1,
		Token:   "TEST_TOKEN",
	})
	updates := []*vk.Update{
		mocks.Update_HandleCommand, mocks.Update_HandleMessage, mocks.Update_HandleMessageNew, mocks.UpdateMessageNew,
		mocks.Update_HandleCommand, mocks.Update_HandleMessage, mocks.Update_HandleMessageNew, mocks.UpdateMessageEdit,
		mocks.UpdateMessageReply, mocks.UpdateMessageAllow, mocks.UpdateMessageDeny,
	}
	m := vk.Middleware{
		OnPreProcessUpdate: func(b *vk.Bot, u *vk.Update) bool {
			OnPreProcessUpdate = true
			return true
		},
		OnPostProcessUpdate: func(b *vk.Bot, u *vk.Update) bool {
			OnPostProcessUpdate = true
			return true
		},
		OnPreMessageNewUpdate: func(b *vk.Bot, u *vk.MessageNew) bool {
			if u.Message.ID != 1 {
				update_rejected = true
				return false
			}
			OnPreMessageNewUpdate = true
			return true
		},
		OnPostMessageNewUpdate: func(b *vk.Bot, u *vk.MessageNew) bool {
			OnPostMessageNewUpdate = true
			return true
		},
		OnPreMessageReplyUpdate: func(b *vk.Bot, u *vk.MessageReply) bool {
			OnPreMessageReplyUpdate = true
			return true
		},
		OnPostMessageReplyUpdate: func(b *vk.Bot, u *vk.MessageReply) bool {
			OnPostMessageReplyUpdate = true
			return true
		},
		OnPreMessageEditUpdate: func(b *vk.Bot, u *vk.MessageEdit) bool {
			OnPreMessageEditUpdate = true
			return true
		},
		OnPostMessageEditUpdate: func(b *vk.Bot, u *vk.MessageEdit) bool {
			OnPostMessageEditUpdate = true
			return true
		},
		OnPreMessageAllowUpdate: func(b *vk.Bot, u *vk.MessageAllow) bool {
			OnPreMessageAllowUpdate = true
			return true
		},
		OnPostMessageAllowUpdate: func(b *vk.Bot, u *vk.MessageAllow) bool {
			OnPostMessageAllowUpdate = true
			return true
		},
		OnPreMessageDenyUpdate: func(b *vk.Bot, u *vk.MessageDeny) bool {
			OnPreMessageDenyUpdate = true
			return true
		},
		OnPostMessageDenyUpdate: func(b *vk.Bot, u *vk.MessageDeny) bool {
			OnPostMessageDenyUpdate = true
			return true
		},
	}

	b.HandleMessageNew(func(b *vk.Bot, m *vk.MessageNew) {})
	b.AddMiddleware(m)
	b.AddMiddleware(vk.LoggingMiddleware)
	b.ProcessUpdates(updates)

	assert.True(t, update_rejected)
	assert.True(t, OnPreProcessUpdate)
	assert.True(t, OnPostProcessUpdate)
	assert.True(t, OnPreMessageNewUpdate)
	assert.True(t, OnPostMessageNewUpdate)
	assert.True(t, OnPreMessageReplyUpdate)
	assert.True(t, OnPostMessageReplyUpdate)
	assert.True(t, OnPreMessageEditUpdate)
	assert.True(t, OnPostMessageEditUpdate)
	assert.True(t, OnPreMessageAllowUpdate)
	assert.True(t, OnPostMessageAllowUpdate)
	assert.True(t, OnPreMessageDenyUpdate)
	assert.True(t, OnPostMessageDenyUpdate)
}
