package service

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
	"gitlab.com/qulaz/khti_timetable_bot/bot/mocks"
	"gitlab.com/qulaz/khti_timetable_bot/bot/parser"
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"log"
	"os"
	"testing"
	"time"
)

func init() {
	if err := os.Chdir(helpers.GetProjectRootDir()); err != nil {
		log.Fatalf("timetable_test: %v", err)
	}
}

func TestNextLesson(t *testing.T) {
	timetable, err := parser.Parse("parser/testdata/timetable.xls")
	assert.NoError(t, err)

	type testCase struct {
		v         time.Time
		groupCode string
		res       string
	}
	testCases := []testCase{
		{
			time.Date(2020, 4, 30, 14, 50, 12, 0, tools.LocalTz),
			"58-1",
			"На сегодня пары закончены",
		},
		{
			time.Date(2020, 4, 30, 14, 00, 12, 0, tools.LocalTz),
			"58-1",
			"4 пара:\nА111 Скуратенко Е.Н. Пр Система государственного и муниципального управления",
		},
		{
			time.Date(2020, 3, 23, 15, 15, 15, 0, tools.LocalTz),
			"58-1",
			"5 пара:\nСпортивный зал №1 Нечаева Л.Н. Пр Консультация по адаптивной физической культуре",
		},
		{
			time.Date(2020, 3, 27, 11, 30, 47, 0, tools.LocalTz),
			"58-1",
			"На сегодня пары закончены",
		},
		{
			time.Date(2020, 3, 24, 11, 40, 57, 0, tools.LocalTz),
			"18-1",
			"3 пара: Окно\n4 пара:\nБ315 Торопов А.С. Лек Теоретические основы электротехники. Часть 2",
		},
		{
			time.Date(2020, 3, 24, 8, 10, 57, 0, tools.LocalTz),
			"58-1",
			"3 пара:\nА229 Зараменских А.А. Лек Мировые информационные ресурсы",
		},
		{
			time.Date(2020, 3, 23, 11, 40, 57, 0, tools.LocalTz),
			"68-1",
			"3 пара: Окно\n4 пара: Окно\n5 пара:\nСпортивный зал №1 Нечаева Л.Н. Пр Консультация по адаптивной физической культуре",
		},
		{
			time.Date(2020, 3, 24, 18, 10, 57, 0, tools.LocalTz),
			"58-1",
			"На сегодня пары закончены",
		},
	}

	for i, tcase := range testCases {
		tools.Now = func() time.Time { return tcase.v }
		res, err := NextLesson(timetable, tcase.groupCode)
		assert.NoErrorf(t, err, "testCase %d", i+1)
		assert.Equalf(t, tcase.res, res, "testCase %d", i+1)
	}
}

func TestTimetableCommand_today_tomorrow(t *testing.T) {
	db.PrepareTestDatabase()
	mocks.InitStartMocks()

	type testCase struct {
		now    time.Time
		peerID int
		body   string
		res    string
	}
	testCases := []testCase{
		{
			// 1 неделя, пятница
			time.Date(2020, 5, 1, 6, 40, 35, 0, tools.LocalTz),
			12, // 58-1
			"сегодня",
			"> Пятница:\n" +
				"1 пара:\nА113 Кадычегова А.Н. Пр Безопасность жизнедеятельности\n" +
				"--------------------\n" +
				"2 пара:\n2-я подгр. А104 Зараменских А.А. Лаб Мировые информационные ресурсы\n" +
				"1-я подгр. А204 Таскин А.Н. Лаб Информационные системы и технологии\n" +
				"--------------------\n" +
				"3 пара:\n-\n" +
				"--------------------\n" +
				"4 пара:\n-\n" +
				"--------------------\n" +
				"5 пара:\n-",
		},
		{
			// 1 неделя, понедельник
			time.Date(2020, 4, 27, 18, 40, 35, 0, tools.LocalTz),
			10, // 19-1
			"сегодня",
			"> Понедельник:\n" +
				"1 пара:\nА223 Чезыбаева Н.В. Пр Иностранный язык\nА224 Танков Е.В. Пр Иностранный язык\n" +
				"--------------------\n" +
				"2 пара:\nА215 Перехожева Е.В. Пр Математический анализ\n" +
				"--------------------\n" +
				"3 пара:\n3-я подгр. А114 (А) Спирин Д.В. Лаб Физика\n" +
				"--------------------\n" +
				"4 пара:\n-\n" +
				"--------------------\n" +
				"5 пара:\nСпортивный зал №1 Нечаева Л.Н. Пр Консультация по адаптивной физической культуре",
		},
		{
			// 1 неделя, воскресенье
			time.Date(2020, 5, 3, 12, 13, 35, 0, tools.LocalTz),
			10, // 19-1
			"сегодня",
			"В воскресенье нет занятий!",
		},
		{
			// 1 неделя, воскресенье
			time.Date(2020, 5, 3, 12, 13, 35, 0, tools.LocalTz),
			10, // 19-1
			"завтра",
			"> Понедельник:\n" +
				"1 пара:\nА223 Чезыбаева Н.В. Пр Иностранный язык\nА224 Танков Е.В. Пр Иностранный язык\n" +
				"--------------------\n" +
				"2 пара:\n-\n" +
				"--------------------\n" +
				"3 пара:\nСпортивный зал №1 Быкова В.А. Пр Физическая культура и спорт\n" +
				"--------------------\n" +
				"4 пара:\n2-я подгр. Б207 Торопов А.С. Лаб Метрология\n" +
				"--------------------\n" +
				"5 пара:\nСпортивный зал №1 Нечаева Л.Н. Пр Консультация по адаптивной физической культуре",
		},
		{
			// 1 неделя, суббота
			time.Date(2020, 5, 2, 12, 13, 35, 0, tools.LocalTz),
			10, // 19-1
			"завтра",
			"В воскресенье нет занятий!",
		},
	}

	for i, tcase := range testCases {
		tools.Now = func() time.Time { return tcase.now }
		mocks.StartMessage.Message.PeerID = tcase.peerID
		mocks.StartMessage.Message.MessageBody = tcase.body
		data := NewData(mocks.StartMessage, common.TimetableCommand, "", MainKeyboard)
		err := TimetableCommand(data)
		assert.NoErrorf(t, err, "testCase %d", i+1)
		assert.Equalf(t, tcase.res, data.Answer, "testCase %d", i+1)
	}
}

func TestTimetableCommand_next_lesson(t *testing.T) {
	db.PrepareTestDatabase()
	mocks.InitStartMocks()
	mocks.StartMessage.Message.MessageBody = "следующая"

	type testCase struct {
		now    time.Time
		peerID int
		res    string
	}
	testCases := []testCase{
		{
			// 1 неделя, четверг
			time.Date(2020, 4, 30, 14, 50, 12, 0, tools.LocalTz),
			12, // 58-1
			"На сегодня пары закончены",
		},
		{
			// 1 неделя, четверг
			time.Date(2020, 4, 30, 14, 00, 12, 0, tools.LocalTz),
			12, // 58-1
			"4 пара:\nА111 Скуратенко Е.Н. Пр Система государственного и муниципального управления",
		},
		{
			// 2 неделя, понедельник
			time.Date(2020, 3, 23, 15, 15, 15, 0, tools.LocalTz),
			12, // 58-1
			"5 пара:\nСпортивный зал №1 Нечаева Л.Н. Пр Консультация по адаптивной физической культуре",
		},
		{
			// 2 неделя, пятница
			time.Date(2020, 3, 27, 11, 30, 47, 0, tools.LocalTz),
			12, // 58-1
			"На сегодня пары закончены",
		},
		{
			// 2 неделя, вторник
			time.Date(2020, 3, 24, 11, 40, 57, 0, tools.LocalTz),
			9, // 18-1
			"3 пара: Окно\n4 пара:\nБ315 Торопов А.С. Лек Теоретические основы электротехники. Часть 2",
		},
		{
			// 2 неделя, вторник
			time.Date(2020, 3, 24, 8, 10, 57, 0, tools.LocalTz),
			12, // 58-1
			"3 пара:\nА229 Зараменских А.А. Лек Мировые информационные ресурсы",
		},
		{
			// 2 неделя, понедельник
			time.Date(2020, 3, 23, 11, 40, 57, 0, tools.LocalTz),
			13, // 68-1
			"3 пара: Окно\n4 пара: Окно\n5 пара:\nСпортивный зал №1 Нечаева Л.Н. Пр Консультация по адаптивной физической культуре",
		},
	}

	for i, tcase := range testCases {
		tools.Now = func() time.Time { return tcase.now }
		mocks.StartMessage.Message.PeerID = tcase.peerID
		data := NewData(mocks.StartMessage, common.TimetableCommand, "", MainKeyboard)
		err := TimetableCommand(data)
		assert.NoErrorf(t, err, "testCase %d", i+1)
		assert.Equalf(t, tcase.res, data.Answer, "testCase %d", i+1)
	}
}

func TestTimetableCommand_default(t *testing.T) {
	db.PrepareTestDatabase()
	mocks.InitStartMocks()

	type testCase struct {
		peerID int
		body   string
		res    string
	}
	testCases := []testCase{
		{
			12, // 58-1
			"1 пятница",
			"> Пятница:\n" +
				"1 пара:\nА113 Кадычегова А.Н. Пр Безопасность жизнедеятельности\n" +
				"--------------------\n" +
				"2 пара:\n2-я подгр. А104 Зараменских А.А. Лаб Мировые информационные ресурсы\n" +
				"1-я подгр. А204 Таскин А.Н. Лаб Информационные системы и технологии\n" +
				"--------------------\n" +
				"3 пара:\n-\n" +
				"--------------------\n" +
				"4 пара:\n-\n" +
				"--------------------\n" +
				"5 пара:\n-",
		},
		{
			13, // 68-1
			"1 понедельник",
			"> Понедельник:\n" +
				"1 пара:\nА105 Сагалакова М.М. Лаб Компьютерная графика\n" +
				"--------------------\n" +
				"2 пара:\nА111 Спирин Д.В. Лек Физика\n" +
				"--------------------\n" +
				"3 пара:\n-\n" +
				"--------------------\n" +
				"4 пара:\n-\n" +
				"--------------------\n" +
				"5 пара:\nСпортивный зал №1 Нечаева Л.Н. Пр Консультация по адаптивной физической культуре",
		},
		{
			10, // 19-1
			"1 понедельник",
			"> Понедельник:\n" +
				"1 пара:\nА223 Чезыбаева Н.В. Пр Иностранный язык\nА224 Танков Е.В. Пр Иностранный язык\n" +
				"--------------------\n" +
				"2 пара:\nА215 Перехожева Е.В. Пр Математический анализ\n" +
				"--------------------\n" +
				"3 пара:\n3-я подгр. А114 (А) Спирин Д.В. Лаб Физика\n" +
				"--------------------\n" +
				"4 пара:\n-\n" +
				"--------------------\n" +
				"5 пара:\nСпортивный зал №1 Нечаева Л.Н. Пр Консультация по адаптивной физической культуре",
		},
		{
			10, // 19-1
			"2 воскресенье",
			"Воскресенье - выходной!",
		},
		{
			10, // 19-1
			"2 понедельник",
			"> Понедельник:\n" +
				"1 пара:\nА223 Чезыбаева Н.В. Пр Иностранный язык\nА224 Танков Е.В. Пр Иностранный язык\n" +
				"--------------------\n" +
				"2 пара:\n-\n" +
				"--------------------\n" +
				"3 пара:\nСпортивный зал №1 Быкова В.А. Пр Физическая культура и спорт\n" +
				"--------------------\n" +
				"4 пара:\n2-я подгр. Б207 Торопов А.С. Лаб Метрология\n" +
				"--------------------\n" +
				"5 пара:\nСпортивный зал №1 Нечаева Л.Н. Пр Консультация по адаптивной физической культуре",
		},
		{
			10, // 19-1
			"0 четверг",
			"Ошибка в команде! Номер недели может принимать значения только 1 или 2",
		},
		{
			10, // 19-1
			"0 субота",
			"Ошибка в команде! Номер недели может принимать значения только 1 или 2",
		},
		{
			10, // 19-1
			"1 субота",
			"Указано неверное название дня недели",
		},
		{
			10, // 19-1
			"1 15 6 субота",
			"Ошибка в команде. Правильное использование команды: \n" +
				"> «/расписание 1 понедельник» - расписание на понедельник первой недели\n" +
				"> «/расписание 2 полное» - полное расписание второй недели",
		},
		{
			10, // 19-1
			"1 15",
			"Указано неверное название дня недели",
		},
		{
			10, // 19-1
			"первая четверг",
			"Ошибка в команде. Правильное использование команды: \n" +
				"> «/расписание 1 понедельник» - расписание на понедельник первой недели\n" +
				"> «/расписание 2 полное» - полное расписание второй недели",
		},
		{
			10, // 19-1
			"1 monday",
			"Указано неверное название дня недели",
		},
	}

	for i, tcase := range testCases {
		mocks.StartMessage.Message.PeerID = tcase.peerID
		mocks.StartMessage.Message.MessageBody = tcase.body
		data := NewData(mocks.StartMessage, common.TimetableCommand, "", MainKeyboard)
		err := TimetableCommand(data)
		assert.NoErrorf(t, err, "testCase %d", i+1)
		assert.Equalf(t, tcase.res, data.Answer, "testCase %d", i+1)
	}
}
