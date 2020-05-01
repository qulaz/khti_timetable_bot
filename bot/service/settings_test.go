package service

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
	"gitlab.com/qulaz/khti_timetable_bot/bot/mocks"
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"testing"
)

func TestSettingsCommand_timetable(t *testing.T) {
	db.PrepareTestDatabase()
	mocks.InitStartMocks()
	mocks.StartMessage.Message.MessageBody = "расписание"
	mocks.StartMessage.Message.PeerID = 2

	data := NewData(mocks.StartMessage, common.SettingsCommand, "", nil)
	err := SettingsCommand(data)
	assert.NoError(t, err)
	assert.Equal(t, "Готово!", data.Answer)
	u, err := db.GetUserByVkID(mocks.StartMessage.Message.PeerID)
	tools.Fatal(t, assert.NoError(t, err))
	assert.False(t, u.IsSubscribed)
	k, ok := data.K.Buttons[1][0].(vk.TextButton)
	tools.Fatal(t, assert.True(t, ok))
	assert.Equal(t, vk.COLOR_POSITIVE, k.Color)

	data = NewData(mocks.StartMessage, common.SettingsCommand, "", nil)
	err = SettingsCommand(data)
	assert.NoError(t, err)
	assert.Equal(t, "Готово!", data.Answer)
	u, err = db.GetUserByVkID(mocks.StartMessage.Message.PeerID)
	tools.Fatal(t, assert.NoError(t, err))
	assert.True(t, u.IsSubscribed)
	k, ok = data.K.Buttons[1][0].(vk.TextButton)
	tools.Fatal(t, assert.True(t, ok))
	assert.Equal(t, vk.COLOR_NEGATIVE, k.Color)
}

func TestSettingsCommand_newsletter(t *testing.T) {
	db.PrepareTestDatabase()
	mocks.InitStartMocks()
	mocks.StartMessage.Message.MessageBody = "рассылка"
	mocks.StartMessage.Message.PeerID = 2

	data := NewData(mocks.StartMessage, common.SettingsCommand, "", nil)
	err := SettingsCommand(data)
	assert.NoError(t, err)
	assert.Equal(t, "Готово!", data.Answer)
	u, err := db.GetUserByVkID(mocks.StartMessage.Message.PeerID)
	tools.Fatal(t, assert.NoError(t, err))
	assert.False(t, u.IsNewsletterEnabled)
	k, ok := data.K.Buttons[2][0].(vk.TextButton)
	tools.Fatal(t, assert.True(t, ok))
	assert.Equal(t, vk.COLOR_POSITIVE, k.Color)

	data = NewData(mocks.StartMessage, common.SettingsCommand, "", nil)
	err = SettingsCommand(data)
	assert.NoError(t, err)
	assert.Equal(t, "Готово!", data.Answer)
	u, err = db.GetUserByVkID(mocks.StartMessage.Message.PeerID)
	tools.Fatal(t, assert.NoError(t, err))
	assert.True(t, u.IsNewsletterEnabled)
	k, ok = data.K.Buttons[2][0].(vk.TextButton)
	tools.Fatal(t, assert.True(t, ok))
	assert.Equal(t, vk.COLOR_NEGATIVE, k.Color)
}

func TestSettingsCommand_ignore(t *testing.T) {
	db.PrepareTestDatabase()
	mocks.InitStartMocks()
	mocks.StartMessage.Message.MessageBody = "unknown"
	mocks.StartMessage.Message.PeerID = 2

	data := NewData(mocks.StartMessage, common.SettingsCommand, "", nil)
	err := SettingsCommand(data)
	assert.Error(t, err)
	assert.Equal(t, common.IgnoreMessageError, err)
}

func TestSettingsCommand_keyboard(t *testing.T) {
	db.PrepareTestDatabase()
	mocks.InitStartMocks()
	mocks.StartMessage.Message.MessageBody = ""
	mocks.StartMessage.Message.PeerID = 2

	u, err := db.GetUserByVkID(mocks.StartMessage.Message.PeerID)
	tools.Fatal(t, assert.NoError(t, err))

	data := NewData(mocks.StartMessage, common.SettingsCommand, "", nil)
	err = SettingsCommand(data)
	tools.Fatal(t, assert.NoError(t, err))
	tcase := "Подробнее о командах настроек:\n" +
		"«Изменить группу» - смена установленной группы. Текущая группа: 17-1;\n" +
		"«Выкл. оповещения о расписании» - выключить оповещения об изменениях в расписании;\n" +
		"«Выкл. новостную рассылку института» - выключить рассылку о важных событиях института и группы;\n"
	assert.Equal(t, tcase, data.Answer)
	subK, ok := data.K.Buttons[1][0].(vk.TextButton)
	tools.Fatal(t, assert.True(t, ok))
	assert.Equal(t, vk.COLOR_NEGATIVE, subK.Color)
	assert.Equal(t, "Выкл. оповещения о расписании", subK.Action.Label)
	newsK, ok := data.K.Buttons[2][0].(vk.TextButton)
	tools.Fatal(t, assert.True(t, ok))
	assert.Equal(t, vk.COLOR_NEGATIVE, newsK.Color)
	assert.Equal(t, "Выкл. новостную рассылку института", newsK.Action.Label)

	tools.Fatal(t, assert.NoError(t, u.SetSubscribe(false)))
	data = NewData(mocks.StartMessage, common.SettingsCommand, "", nil)
	err = SettingsCommand(data)
	tools.Fatal(t, assert.NoError(t, err))
	tcase = "Подробнее о командах настроек:\n" +
		"«Изменить группу» - смена установленной группы. Текущая группа: 17-1;\n" +
		"«Вкл. оповещения о расписании» - включить оповещения об изменениях в расписании;\n" +
		"«Выкл. новостную рассылку института» - выключить рассылку о важных событиях института и группы;\n"
	assert.Equal(t, tcase, data.Answer)
	subK, ok = data.K.Buttons[1][0].(vk.TextButton)
	tools.Fatal(t, assert.True(t, ok))
	assert.Equal(t, vk.COLOR_POSITIVE, subK.Color)
	assert.Equal(t, "Вкл. оповещения о расписании", subK.Action.Label)
	newsK, ok = data.K.Buttons[2][0].(vk.TextButton)
	tools.Fatal(t, assert.True(t, ok))
	assert.Equal(t, vk.COLOR_NEGATIVE, newsK.Color)
	assert.Equal(t, "Выкл. новостную рассылку института", newsK.Action.Label)

	tools.Fatal(t, assert.NoError(t, u.SetNewsletterEnabling(false)))
	data = NewData(mocks.StartMessage, common.SettingsCommand, "", nil)
	err = SettingsCommand(data)
	tools.Fatal(t, assert.NoError(t, err))
	tcase = "Подробнее о командах настроек:\n" +
		"«Изменить группу» - смена установленной группы. Текущая группа: 17-1;\n" +
		"«Вкл. оповещения о расписании» - включить оповещения об изменениях в расписании;\n" +
		"«Вкл. новостную рассылку института» - включить рассылку о важных событиях института и группы;\n"
	assert.Equal(t, tcase, data.Answer)
	subK, ok = data.K.Buttons[1][0].(vk.TextButton)
	tools.Fatal(t, assert.True(t, ok))
	assert.Equal(t, vk.COLOR_POSITIVE, subK.Color)
	assert.Equal(t, "Вкл. оповещения о расписании", subK.Action.Label)
	newsK, ok = data.K.Buttons[2][0].(vk.TextButton)
	tools.Fatal(t, assert.True(t, ok))
	assert.Equal(t, vk.COLOR_POSITIVE, newsK.Color)
	assert.Equal(t, "Вкл. новостную рассылку института", newsK.Action.Label)
}
