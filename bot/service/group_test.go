package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
	"gitlab.com/qulaz/khti_timetable_bot/bot/mocks"
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"gopkg.in/khaiql/dbcleaner.v2"
	"gopkg.in/khaiql/dbcleaner.v2/engine"
	"testing"
)

var Cleaner = dbcleaner.New()

type ServiceTestSuite struct {
	suite.Suite
}

func (suite *ServiceTestSuite) SetupSuite() {
	common.TestInits()
	db.InitTestDatabase()
	pg := engine.NewPostgresEngine(helpers.Config.POSTGRES_DSN)
	Cleaner.SetEngine(pg)
}

func (suite *ServiceTestSuite) TearDownSuite() {
	Cleaner.Clean("groups", "timetable", "users")
	db.CloseDatabase()
}

func (suite *ServiceTestSuite) SetupTest() {
	Cleaner.Acquire("groups", "timetable", "users")
	db.PrepareTestDatabase()
}

func (suite *ServiceTestSuite) Test_buildGroupKeyboard_success() {
	t := suite.T()

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

func (suite *ServiceTestSuite) Test_buildGroupKeyboard_tooBigLimit() {
	t := suite.T()

	k, err := buildGroupKeyboard(groupsLimit+1, 0)
	assert.Error(t, err)
	assert.Nil(t, k)
}

func (suite *ServiceTestSuite) TestGroupCommand_success_pagination() {
	t := suite.T()

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

func (suite *ServiceTestSuite) TestGroupCommand_nil_payload() {
	t := suite.T()

	mocks.InitStartMocks()
	mocks.StartMessage.Message.MessageCommand = "/group"
	mocks.StartMessage.Message.MessageBody = nextBody
	mocks.StartMessage.Message.Payload = nil

	data := NewData(mocks.StartMessage, common.GroupCommand, common.UnknownErrorMessage, vk.NewKeyboard(true))
	err := GroupCommand(data)
	assert.Error(t, err)
}

func (suite *ServiceTestSuite) TestGroupCommand_unknown_group() {
	t := suite.T()

	mocks.InitStartMocks()
	mocks.StartMessage.Message.MessageBody = "55-5"

	data := NewData(mocks.StartMessage, common.GroupCommand, common.UnknownErrorMessage, vk.NewKeyboard(true))
	err := GroupCommand(data)
	assert.Error(t, err)
}

func (suite *ServiceTestSuite) TestGroupCommand_success_new_and_change() {
	t := suite.T()

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

func (suite *ServiceTestSuite) TestGroupCommand_nonkeyboard_answer() {
	t := suite.T()

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

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
