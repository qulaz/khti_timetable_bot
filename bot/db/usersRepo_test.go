package db

import (
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
)

const UserFixtureCount = 13

func (suite *DBTestSuite) TestCreateUser() {
	vkid := 50

	err := CreateUser(vkid, 1)
	tools.Fatal(suite.T(), suite.NoError(err))

	u, err := GetUserByVkID(vkid)
	tools.Fatal(suite.T(), suite.NoError(err))
	suite.Equal(vkid, u.VkID)
	suite.Falsef(!u.IsNewsletterEnabled || !u.IsSubscribed || !u.IsActive, "user: %+v", u)

	err = CreateUser(1, 5)
	suite.Error(err)
}

func (suite *DBTestSuite) TestGetUserByVkID() {
	vkid := 1

	u, err := GetUserByVkID(vkid)
	tools.Fatal(suite.T(), suite.NoError(err))
	suite.Equal(vkid, u.VkID)
	suite.Falsef(!u.IsNewsletterEnabled || !u.IsSubscribed || !u.IsActive, "user: %+v", u)
}

func (suite *DBTestSuite) TestGetUsers() {
	type testCase struct {
		limit  int
		offset int
		len    int
	}
	testCases := []testCase{
		{100, 0, UserFixtureCount},
		{1, 0, 1},
		{5, 0, 5},
		{100, 10, UserFixtureCount - 10},
		{100, 1, UserFixtureCount - 1},
		{100, UserFixtureCount, 0},
		{100, 999, 0},
	}

	for i, tcase := range testCases {
		users, err := GetUsers(tcase.limit, tcase.offset)
		tools.Fatal(suite.T(), suite.NoErrorf(err, "testCase: %d", i+1))
		suite.Equal(tcase.len, len(users))
	}
}

func (suite *DBTestSuite) TestUserModel_SetSubscribe() {
	type testCase struct {
		vkid int
		v    bool
	}
	testCases := []testCase{
		{1, false},
		{1, true},
		{1, true},
		{6, false},
		{6, true},
		{7, true},
	}

	for i, tcase := range testCases {
		u, err := GetUserByVkID(tcase.vkid)
		tools.Fatal(suite.T(), suite.NoErrorf(err, "testCase: %d", i+1))

		err = u.SetSubscribe(tcase.v)
		tools.Fatal(suite.T(), suite.NoErrorf(err, "testCase: %d", i+1))
		suite.Equal(tcase.v, u.IsSubscribed)
	}
}

func (suite *DBTestSuite) TestUserModel_SetNewsletterEnabling() {
	type testCase struct {
		vkid int
		v    bool
	}
	testCases := []testCase{
		{1, false},
		{1, true},
		{1, true},
		{6, false},
		{6, true},
		{7, true},
	}

	for i, tcase := range testCases {
		u, err := GetUserByVkID(tcase.vkid)
		tools.Fatal(suite.T(), suite.NoErrorf(err, "testCase: %d", i+1))

		err = u.SetNewsletterEnabling(tcase.v)
		tools.Fatal(suite.T(), suite.NoErrorf(err, "testCase: %d", i+1))
		suite.Equal(tcase.v, u.IsNewsletterEnabled)
	}
}

func (suite *DBTestSuite) TestUserModel_SetActive() {
	type testCase struct {
		vkid int
		v    bool
	}
	testCases := []testCase{
		{1, false},
		{1, true},
		{1, true},
		{5, true},
		{6, false},
		{7, true},
		{8, false},
	}

	for i, tcase := range testCases {
		u, err := GetUserByVkID(tcase.vkid)
		tools.Fatal(suite.T(), suite.NoErrorf(err, "testCase: %d", i+1))

		err = u.SetActive(tcase.v)
		tools.Fatal(suite.T(), suite.NoErrorf(err, "testCase: %d", i+1))
		suite.Equal(tcase.v, u.IsActive)
	}
}

func (suite *DBTestSuite) TestUserModel_ChangeGroup() {
	type testCase struct {
		vkid int
		v    string
	}
	testCases := []testCase{
		{1, "58-1"},
		{1, "58-1"},
		{2, "19-1"},
		{3, "18-1"},
	}

	for i, tcase := range testCases {
		u, err := GetUserByVkID(tcase.vkid)
		tools.Fatal(suite.T(), suite.NoErrorf(err, "testCase: %d", i+1))

		err = u.ChangeGroup(tcase.v)
		tools.Fatal(suite.T(), suite.NoErrorf(err, "testCase: %d", i+1))
		suite.Equal(tcase.v, u.Group.Code)
	}

	u, err := GetUserByVkID(5)
	tools.Fatal(suite.T(), suite.NoError(err))
	code := u.Group.Code

	err = u.ChangeGroup("99-5")
	tools.Fatal(suite.T(), suite.Error(err))
	suite.Equal(code, u.Group.Code)
}
