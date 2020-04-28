package service

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
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
