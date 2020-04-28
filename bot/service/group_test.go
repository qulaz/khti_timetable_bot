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

func TestMain(m *testing.M) {
	db.TestMainWithDb(m)
}

func Test_buildGroupKeyboard_success(t *testing.T) {
	db.PrepareTestDatabase()

	total, err := db.GroupsCount()
	tools.Fatal(t, assert.NoError(t, err))

	k, err := buildGroupKeyboard(groupsLimit, 0)
	tools.Fatal(t, assert.NoError(t, err))
	assert.Equal(t, groupsLimit+1, k.ButtonCount())
	assert.Equal(t, 10, k.RowsCount())
	assert.Equal(t, 27, k.Buttons[9][0].(vk.TextButton).Action.Payload.Offset)

	k, err = buildGroupKeyboard(groupsLimit, 27)
	tools.Fatal(t, assert.NoError(t, err))
	assert.Equal(t, total-groupsLimit+1, k.ButtonCount())
	assert.Equal(t, 0, k.Buttons[len(k.Buttons)-1][0].(vk.TextButton).Action.Payload.Offset)
}

func Test_buildGroupKeyboard_tooBigLimit(t *testing.T) {
	k, err := buildGroupKeyboard(groupsLimit+1, 0)
	assert.Error(t, err)
	assert.Nil(t, k)
}

func TestGroupCommand_success_pagination(t *testing.T) {
	db.PrepareTestDatabase()
	mocks.InitStartMocks()
	mocks.StartMessage.Message.MessageCommand = "/group"
	mocks.StartMessage.Message.MessageBody = nextBody
	mocks.StartMessage.Message.Payload = &vk.ButtonPayload{Command: "/group", Body: nextBody, Offset: 27}

	total, err := db.GroupsCount()
	tools.Fatal(t, assert.NoError(t, err))

	data := NewData(mocks.StartMessage, common.GroupCommand, common.UnknownErrorMessage, vk.NewKeyboard(true))
	err = GroupCommand(data)
	tools.Fatal(t, assert.NoError(t, err))
	assert.Equal(t, "Выбери группу, в которой ты учишься", data.Answer)
	assert.Equal(t, total-27+1, data.K.ButtonCount())

	mocks.StartMessage.Message.MessageBody = prevBody
	mocks.StartMessage.Message.Payload = &vk.ButtonPayload{Command: "/group", Body: prevBody, Offset: 0}
	data = NewData(mocks.StartMessage, common.GroupCommand, common.UnknownErrorMessage, vk.NewKeyboard(true))
	err = GroupCommand(data)
	tools.Fatal(t, assert.NoError(t, err))
	assert.Equal(t, "Выбери группу, в которой ты учишься", data.Answer)
	assert.Equal(t, groupsLimit+1, data.K.ButtonCount())
}

func TestGroupCommand_nil_payload(t *testing.T) {
	db.PrepareTestDatabase()
	mocks.InitStartMocks()
	mocks.StartMessage.Message.MessageCommand = "/group"
	mocks.StartMessage.Message.MessageBody = nextBody
	mocks.StartMessage.Message.Payload = nil

	data := NewData(mocks.StartMessage, common.GroupCommand, common.UnknownErrorMessage, vk.NewKeyboard(true))
	err := GroupCommand(data)
	assert.Error(t, err)
}

func TestGroupCommand_unknown_group(t *testing.T) {
	db.PrepareTestDatabase()
	mocks.InitStartMocks()
	mocks.StartMessage.Message.MessageBody = "55-5"

	data := NewData(mocks.StartMessage, common.GroupCommand, common.UnknownErrorMessage, vk.NewKeyboard(true))
	err := GroupCommand(data)
	assert.Error(t, err)
}

func TestGroupCommand_success_new_and_change(t *testing.T) {
	db.PrepareTestDatabase()
	mocks.InitStartMocks()
	mocks.StartMessage.Message.MessageBody = "58-1"

	data := NewData(mocks.StartMessage, common.GroupCommand, common.UnknownErrorMessage, vk.NewKeyboard(true))
	err := GroupCommand(data)
	tools.Fatal(t, assert.NoError(t, err))
	u, err := db.GetUserByVkID(mocks.StartMessage.Message.PeerID)
	tools.Fatal(t, assert.NoError(t, err))
	assert.Equal(t, "58-1", u.Group.Code)
	assert.Contains(t, data.Answer, "Группа выбрана!")
	assert.Contains(t, data.Answer, MainKeyboardHelp)

	mocks.StartMessage.Message.MessageBody = "57-1"
	data = NewData(mocks.StartMessage, common.GroupCommand, common.UnknownErrorMessage, vk.NewKeyboard(true))
	err = GroupCommand(data)
	tools.Fatal(t, assert.NoError(t, err))
	u, err = db.GetUserByVkID(mocks.StartMessage.Message.PeerID)
	tools.Fatal(t, assert.NoError(t, err))
	assert.Equal(t, "57-1", u.Group.Code)
	assert.Equal(t, "Группа успешно изменена!", data.Answer)
}

func TestGroupCommand_nonkeyboard_answer(t *testing.T) {
	db.PrepareTestDatabase()
	mocks.InitStartMocks()
	mocks.StartMessage.Message.MessageBody = "58-1"
	mocks.StartMessage.ClientInfo.Keyboard = false

	data := NewData(mocks.StartMessage, common.GroupCommand, common.UnknownErrorMessage, vk.NewKeyboard(true))
	err := GroupCommand(data)
	tools.Fatal(t, assert.NoError(t, err))
	u, err := db.GetUserByVkID(mocks.StartMessage.Message.PeerID)
	tools.Fatal(t, assert.NoError(t, err))
	assert.Equal(t, "58-1", u.Group.Code)
	assert.Contains(t, data.Answer, "Группа выбрана!")
	assert.Contains(t, data.Answer, NonKeyboardMainHelp)
}
