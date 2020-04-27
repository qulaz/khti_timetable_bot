package db

import "testing"

const UserFixtureCount = 11

func TestCreateUser(t *testing.T) {
	PrepareTestDatabase()

	vkid := 50
	if err := CreateUser(vkid, 1); err != nil {
		t.Fatalf("Ошибка создания пользователя: %+v\n", err)
	}
	u, err := GetUserByVkID(vkid)
	if err != nil {
		t.Fatalf("Созданный пользователь не найден: %+v\n", err)
	}
	if u.VkID != vkid {
		t.Errorf("vk_id созданного пользователя не совпадает с заданным %d != %d\n", u.VkID, vkid)
	}
	if !u.IsNewsletterEnabled || !u.IsSubscribed || !u.IsActive {
		t.Errorf("Стандартные значения не совпадают с ожидаемыми: %+v\n", u)
	}

	if err := CreateUser(1, 5); err == nil {
		t.Error("Создан пользователь с дублирующимся vk_id")
	}
}

func TestGetUserByVkID(t *testing.T) {
	PrepareTestDatabase()

	vkid := 1
	u, err := GetUserByVkID(vkid)
	if err != nil {
		t.Fatalf("Ошибка получения пользователя: %+v\n", err)
	}
	if u.VkID != vkid {
		t.Errorf("vk_id созданного пользователя не совпадает с заданным %d != %d\n", u.VkID, vkid)
	}
	if !u.IsNewsletterEnabled || !u.IsSubscribed || !u.IsActive {
		t.Errorf("Стандартные значения не совпадают с ожидаемыми: %+v\n", u)
	}
}

func TestGetUsers(t *testing.T) {
	PrepareTestDatabase()

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
		if err != nil {
			t.Errorf("%d. Ошибка при получении списка пользователей в тест кейсе: %+v\n", i+1, err)
			continue
		}
		if lenG := len(users); lenG != tcase.len {
			t.Errorf("%d. Ожидаемое кол-во пользователей %d не равно действительному %d\n", i+1, tcase.len, lenG)
		}
	}
}

func TestUserModel_SetSubscribe(t *testing.T) {
	PrepareTestDatabase()

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
		if err != nil {
			t.Errorf("%d. Ошибка получения пользователя: %+v\n", i+1, err)
			continue
		}
		if err := u.SetSubscribe(tcase.v); err != nil {
			t.Errorf("%d. Ошибка смены подписки пользователя: %+v\n", i+1, err)
			continue
		}
		if u.IsSubscribed != tcase.v {
			t.Errorf("%d. Ошибка u.IsSubscribed != tcase.v; %v != %v\n", i+1, u.IsSubscribed, tcase.v)
		}
	}
}

func TestUserModel_SetNewsletterEnabling(t *testing.T) {
	PrepareTestDatabase()

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
		if err != nil {
			t.Errorf("%d. Ошибка получения пользователя: %+v\n", i+1, err)
			continue
		}
		if err := u.SetNewsletterEnabling(tcase.v); err != nil {
			t.Errorf("%d. Ошибка смены подписки пользователя: %+v\n", i+1, err)
			continue
		}
		if u.IsNewsletterEnabled != tcase.v {
			t.Errorf("%d. Ошибка u.IsNewsletterEnabled != tcase.v; %v != %v\n", i+1, u.IsNewsletterEnabled, tcase.v)
		}
	}
}

func TestUserModel_SetActive(t *testing.T) {
	PrepareTestDatabase()

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
		if err != nil {
			t.Errorf("%d. Ошибка получения пользователя: %+v\n", i+1, err)
			continue
		}
		if err := u.SetActive(tcase.v); err != nil {
			t.Errorf("%d. Ошибка смены типа активности пользователя: %+v\n", i+1, err)
			continue
		}
		if u.IsActive != tcase.v {
			t.Errorf("%d. Ошибка u.IsActive != tcase.v; %v != %v\n", i+1, u.IsActive, tcase.v)
		}
	}
}

func TestUserModel_ChangeGroup(t *testing.T) {
	PrepareTestDatabase()

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
		if err != nil {
			t.Errorf("%d. Ошибка получения пользователя: %+v\n", i+1, err)
			continue
		}
		if err := u.ChangeGroup(tcase.v); err != nil {
			t.Errorf("%d. Ошибка смены группы пользователя: %+v\n", i+1, err)
			continue
		}
		if u.Group.Code != tcase.v {
			t.Errorf("%d. Ошибка u.Group.Code != tcase.v; %s != %s\n", i+1, u.Group.Code, tcase.v)
		}
	}

	u, err := GetUserByVkID(5)
	code := u.Group.Code
	if err != nil {
		t.Errorf("Ошибка получения пользователя: %+v\n", err)
	}
	if err := u.ChangeGroup("99-5"); err == nil {
		t.Fatalf("Группа сменена на несуществующую\n")
	}
	if u.Group.Code != code {
		t.Errorf("Поменялась группа! u.Group.Code != code; %s != %s\n", u.Group.Code, code)
	}
}
