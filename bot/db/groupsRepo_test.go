package db

import (
	"testing"
)

const GroupFixtureCount = 33

func TestMain(m *testing.M) {
	TestMainWithDb(m)
}

func TestGroupsCount(t *testing.T) {
	PrepareTestDatabase()

	c, err := GroupsCount()
	if err != nil {
		t.Errorf("%+v\n", err)
	}

	if c != GroupFixtureCount {
		t.Errorf("Ожидаемое кол-во групп %d не равно действительному %d\n", GroupFixtureCount, c)
	}
}

func TestCreateGroupIfNotExists(t *testing.T) {
	PrepareTestDatabase()

	code := "99-1"
	if err := CreateGroupIfNotExists(code); err != nil {
		t.Fatalf("Ошибка при создании группы: %+v\n", err)
	}

	g, err := GetGroupByGroupCode(code)
	if err != nil {
		t.Fatalf("Ошибка получении группы: %+v\n", err)
	}

	if g.Code != code {
		t.Fatalf("Ожидаемый код группы %q не совпадает с действительным %q\n", code, g.Code)
	}
}

func TestGetGroups(t *testing.T) {
	PrepareTestDatabase()

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
			t.Errorf("%d. Ошибка при получении списка групп в тест кейсе: %+v\n", i+1, err)
			continue
		}
		if lenG := len(groups); lenG != tcase.len {
			t.Errorf("%d. Ожидаемое кол-во групп %d не равно действительному %d\n", i+1, tcase.len, lenG)
		}
	}
}

func TestGetGroup(t *testing.T) {
	PrepareTestDatabase()

	if g, err := GetGroup(1); err != nil {
		t.Fatalf("Ошибка при получении группы с ID 1: %+v\n", err)
	} else {
		if g.Code != "16-1" {
			t.Fatalf("Ожидаемый код группы %q не совпадает с действительным %q\n", "16-1", g.Code)
		}
	}
}

func TestGetGroupByGroupCode(t *testing.T) {
	PrepareTestDatabase()

	code := "58-1"
	if g, err := GetGroupByGroupCode(code); err != nil {
		t.Fatalf("Ошибка при получении группы %q: %+v\n", code, err)
	} else {
		if g.Code != code {
			t.Fatalf("Ожидаемый код группы %q не совпадает с действительным %q\n", code, g.Code)
		}
	}
}

func TestCreateGroupWithTransaction(t *testing.T) {
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Ошибка создания транзакции: %+v\n", err)
	}

	codes := []string{"99-1", "98-1", "97-1", "96-1", "99-2", "99-3", "98-2", "98-3", "97-2", "97-3", "96-2", "96-3"}
	codesCount := len(codes)
	for i, code := range codes {
		if i == codesCount-1 {
			if err := tx.Rollback(); err != nil {
				t.Fatalf("Ошибка отмены транзакции: %+v\n", err)
			}
			break
		}

		if err := createGroupIfNotExists(tx, code); err != nil {
			t.Fatalf("Ошибка при создании группы с кодом %q: %+v\n", code, err)
		}
	}

	for _, code := range codes {
		if _, err := GetGroupByGroupCode(code); err == nil {
			t.Errorf("Группа с кодом %q создалась не смотря на отмену транзакции\n", code)
		}
	}
}
