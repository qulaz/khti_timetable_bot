package service

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/mocks"
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
)

func (suite *ServiceTestSuite) TestStartCommand_success() {
	t := suite.T()

	mocks.InitStartMocks()

	data := NewData(mocks.StartMessage, common.StartCommand, common.UnknownErrorMessage, vk.NewKeyboard(true))
	err := StartCommand(data)
	tools.Fatal(t, assert.NoError(t, err))
	assert.Equal(t, startMessage, data.Answer)
	tools.Fatal(t, assert.NotNil(t, data.K))
	assert.Equal(t, groupsLimit+1, data.K.ButtonCount())
}

func (suite *ServiceTestSuite) TestStartCommand_reset() {
	t := suite.T()

	mocks.InitStartMocks()
	mocks.StartMessage.Message.PeerID = 1
	mocks.StartMessage.Message.MessageBody = "reset"
	mocks.StartMessage.Message.MessageCommand = "/start"

	data := NewData(mocks.StartMessage, common.StartCommand, common.UnknownErrorMessage, vk.NewKeyboard(true))
	err := StartCommand(data)
	tools.Fatal(t, assert.NoError(t, err))
	assert.Equal(t, "Выбери группу, в которой ты учишься", data.Answer)
	tools.Fatal(t, assert.NotNil(t, data.K))
	assert.Equal(t, groupsLimit+1, data.K.ButtonCount())
}

func (suite *ServiceTestSuite) TestStartCommand_keyboard_unsupported() {
	t := suite.T()

	mocks.InitStartMocks()
	mocks.StartMessage.ClientInfo.Keyboard = false

	data := NewData(mocks.StartMessage, common.StartCommand, common.UnknownErrorMessage, vk.NewKeyboard(true))
	err := StartCommand(data)
	tools.Fatal(t, assert.NoError(t, err))
	assert.True(t, assert.NotEqual(t, startMessage, data.Answer) && assert.Contains(t, data.Answer, startMessage))
}

func (suite *ServiceTestSuite) TestStartCommand_already_registered_user() {
	t := suite.T()

	mocks.InitStartMocks()
	mocks.StartMessage.Message.PeerID = 1

	data := NewData(mocks.StartMessage, common.StartCommand, common.UnknownErrorMessage, vk.NewKeyboard(true))
	err := StartCommand(data)
	tools.Fatal(t, assert.Error(t, err))
	assert.Equal(t, RegisteredUserTryingToUseStartCommandError, err)
}
