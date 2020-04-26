package vk_test

import (
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"log"
)

func Example() {
	b, _ := vk.CreateBot(vk.Settings{
		GroupID: 123,
		Token:   "group_token",
	})

	b.HandleCommand("/ping", func(b *vk.Bot, m *vk.MessageNew) {
		mID, err := b.SendTextMessage("pong", m.Message.PeerID)
		if err != nil {
			b.Logger.Error(err)
		}
		b.Logger.Debugf("ID отправленного сообщения: %d", mID)
	})

	if err := b.Run(); err != nil {
		log.Fatal(err)
	}
}

func ExampleBot_AddMiddleware() {
	b, _ := vk.CreateBot(vk.Settings{
		GroupID: 123,
		Token:   "group_token",
	})

	b.AddMiddleware(vk.LoggingMiddleware)
	b.AddMiddleware(vk.Middleware{
		OnPreMessageNewUpdate: func(b *vk.Bot, u *vk.MessageNew) bool {
			AdminID := 1
			if u.Message.MessageCommand == "/ban" && u.Message.FromID != AdminID {
				return false
			}
			return true
		},
	})
}
