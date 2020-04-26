package vk

import (
	"testing"
)

func TestCreateBot(t *testing.T) {
	const (
		TEST_GROUP_ID = 15
		TEST_TOKEN    = "token"
	)

	s := Settings{}
	_, err := CreateBot(s)
	if err == nil {
		t.Error("Бот создался с пустыми настройками")
	}

	s = Settings{GroupID: TEST_GROUP_ID}
	_, err = CreateBot(s)
	if err == nil {
		t.Error("Бот создался без токена")
	}

	s = Settings{Token: TEST_TOKEN}
	_, err = CreateBot(s)
	if err == nil {
		t.Error("Бот создался без группы")
	}

	s = Settings{BaseURL: "random_string"}
	_, err = CreateBot(s)
	if err == nil {
		t.Error("Бот создался c неправильным базовым URL")
	}

	s = Settings{Token: TEST_TOKEN, GroupID: TEST_GROUP_ID}
	b, err := CreateBot(s)
	if err != nil {
		t.Fatal("Бот не создался со всеми обязательными параметрами")
	}

	if b.baseURL.String() != VK_API_BASE_URL {
		t.Errorf("%q != %q", b.baseURL, VK_API_BASE_URL)
	}
	if b.groupID != TEST_GROUP_ID {
		t.Errorf("%d != %d", b.groupID, TEST_GROUP_ID)
	}
	if b.token != TEST_TOKEN {
		t.Errorf("%q != %q", b.token, TEST_TOKEN)
	}
	if b.client == nil {
		t.Error("Не создался клиент")
	}
	if b.Logger == nil {
		t.Error("Не создался логгер")
	}
	if b.poller == nil {
		t.Error("Не создался poller")
	}

	b, err = CreateBot(Settings{
		GroupID: TEST_GROUP_ID,
		Token:   TEST_TOKEN,
		Poller:  NewDefaultLongPoller(),
		Logger:  &DefaultLogType{},
		Client:  NewClient(),
		BaseURL: VK_API_BASE_URL,
	})
	if err != nil {
		t.Fatal("Бот не создался со всеми обязательными параметрами")
	}

	if b.baseURL.String() != VK_API_BASE_URL {
		t.Errorf("%q != %q", b.baseURL.String(), VK_API_BASE_URL)
	}
	if b.groupID != TEST_GROUP_ID {
		t.Errorf("%d != %d", b.groupID, TEST_GROUP_ID)
	}
	if b.token != TEST_TOKEN {
		t.Errorf("%q != %q", b.token, TEST_TOKEN)
	}
	if b.client == nil {
		t.Error("Не создался клиент")
	}
	if b.Logger == nil {
		t.Error("Не создался логгер")
	}
	if b.poller == nil {
		t.Error("Не создался poller")
	}
}

func TestBot_AddMiddleware(t *testing.T) {
	b, _ := CreateBot(Settings{
		GroupID: 15,
		Token:   "ete",
	})

	m := Middleware{
		OnPreProcessUpdate:       func(b *Bot, u *Update) bool { return true },
		OnPostProcessUpdate:      func(b *Bot, u *Update) bool { return true },
		OnPreMessageNewUpdate:    func(b *Bot, u *MessageNew) bool { return true },
		OnPostMessageNewUpdate:   func(b *Bot, u *MessageNew) bool { return true },
		OnPreMessageReplyUpdate:  func(b *Bot, u *MessageReply) bool { return true },
		OnPostMessageReplyUpdate: func(b *Bot, u *MessageReply) bool { return true },
		OnPreMessageEditUpdate:   func(b *Bot, u *MessageEdit) bool { return true },
		OnPostMessageEditUpdate:  func(b *Bot, u *MessageEdit) bool { return true },
		OnPreMessageAllowUpdate:  func(b *Bot, u *MessageAllow) bool { return true },
		OnPostMessageAllowUpdate: func(b *Bot, u *MessageAllow) bool { return true },
		OnPreMessageDenyUpdate:   func(b *Bot, u *MessageDeny) bool { return true },
		OnPostMessageDenyUpdate:  func(b *Bot, u *MessageDeny) bool { return true },
	}

	b.AddMiddleware(LoggingMiddleware)
	b.AddMiddleware(m)

	for key, m := range b.middlewares {
		if l := len(m); l != 2 {
			t.Errorf("Middleware %s имеет длину %d", key, l)
		}
	}
}
