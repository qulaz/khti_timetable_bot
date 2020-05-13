package db

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"gopkg.in/khaiql/dbcleaner.v2"
	"gopkg.in/khaiql/dbcleaner.v2/engine"
	"testing"
)

var Cleaner = dbcleaner.New()

const GroupFixtureCount = 33

type DBTestSuite struct {
	suite.Suite
}

func (suite *DBTestSuite) SetupSuite() {
	common.TestInits()
	InitTestDatabase()
	db := engine.NewPostgresEngine(helpers.Config.POSTGRES_DSN)
	Cleaner.SetEngine(db)
}

func (suite *DBTestSuite) TearDownSuite() {
	Cleaner.Clean("groups", "timetable", "users")
	CloseDatabase()
}

func (suite *DBTestSuite) SetupTest() {
	Cleaner.Acquire("groups", "timetable", "users")
	PrepareTestDatabase()
}

func (suite *DBTestSuite) TestGroupsCount() {
	c, err := GroupsCount()
	tools.Fatal(suite.T(), assert.NoError(suite.T(), err))
	assert.Equal(suite.T(), GroupFixtureCount, c)
}

func (suite *DBTestSuite) TestCreateGroupIfNotExists() {
	code := "99-1"
	tools.Fatal(suite.T(), suite.NoError(CreateGroupIfNotExists(code)))

	g, err := GetGroupByGroupCode(code)
	tools.Fatal(suite.T(), suite.NoError(err))
	suite.Equal(code, g.Code)
}

func (suite *DBTestSuite) TestGetGroups() {
	type testCase struct {
		limit  int
		offset int
		len    int
	}
	testCases := []testCase{
		{100, 0, GroupFixtureCount},
		{1, 0, 1},
		{5, 0, 5},
		{100, 10, GroupFixtureCount - 10},
		{100, GroupFixtureCount, 0},
		{100, 999, 0},
	}

	for i, tcase := range testCases {
		groups, err := GetGroups(tcase.limit, tcase.offset)
		if err != nil {
			suite.T().Errorf("%d. Ошибка при получении списка групп в тест кейсе: %+v\n", i+1, err)
			continue
		}
		suite.Equalf(tcase.len, len(groups), "testCase %d", i+1)
	}
}

func (suite *DBTestSuite) TestGetGroup() {
	g, err := GetGroup(1)
	tools.Fatal(suite.T(), suite.NoError(err))
	suite.Equal("16-1", g.Code)
}

func (suite *DBTestSuite) TestGetGroupByGroupCode() {
	code := "58-1"
	g, err := GetGroupByGroupCode(code)
	tools.Fatal(suite.T(), suite.NoError(err))
	suite.Equal(code, g.Code)
}

func (suite *DBTestSuite) TestCreateGroupWithTransaction() {
	tx, err := db.Begin()
	tools.Fatal(suite.T(), suite.NoError(err))

	codes := []string{"99-1", "98-1", "97-1", "96-1", "99-2", "99-3", "98-2", "98-3", "97-2", "97-3", "96-2", "96-3"}
	codesCount := len(codes)
	for i, code := range codes {
		if i == codesCount-1 {
			err := tx.Rollback()
			tools.Fatal(suite.T(), suite.NoError(err))
			break
		}

		err := createGroupIfNotExists(tx, code)
		tools.Fatal(suite.T(), suite.NoError(err))
	}

	for _, code := range codes {
		_, err := GetGroupByGroupCode(code)
		suite.Errorf(err, "code: %s", code)
	}
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
