package service

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"strconv"
	"strings"
	"time"
)

// Определение следующего занятия
func NextLesson(t *db.Timetable, groupCode string) (string, error) {
	var answer string

	n := tools.Now()
	currentLessonNum := CurrentLessonNum()

	if currentLessonNum == 999 || currentLessonNum+1 > 5 {
		return "На сегодня пары закончены", nil
	}
	if currentLessonNum == 1000 {
		helpers.Logger.Errorw("currentLessonNum==1000", "now", n)
		return "", errors.New("Неизвестная ошибка")
	}

	weekSchedule, err := t.GetWeekSchedule(groupCode, tools.GetCurrentWeekNum())
	if err != nil {
		return "", errors.WithStack(err)
	}

	daySchedule, err := weekSchedule.GetDaySchedule(tools.TodayName())
	if err != nil {
		return "", errors.WithStack(err)
	}
	daySchedule = tools.WindowAnalyzer(daySchedule)

	for nextLesson := currentLessonNum; nextLesson < len(daySchedule); nextLesson++ {
		if nextLesson == -1 {
			continue
		}

		lesson := daySchedule[nextLesson]
		if lesson == "-" {
			answer += fmt.Sprintf("%d пара: Окно\n", nextLesson+1)
		}
		if lesson != "" && lesson != "-" {
			answer += fmt.Sprintf("%d пара:\n%s", nextLesson+1, lesson)
			return answer, nil
		}
	}

	return "На сегодня пары закончены", nil
}

// Парсинг тела команды /расписание. Тело должно иметь вид <weekNum> <dayName>. Прим: 1 понедельник; 2 четверг; 1 полное
func parseBody(body string) (int, string, error) {
	var (
		wrongBody = "Ошибка в команде. Правильное использование команды: \n" +
			"> «/расписание 1 понедельник» - расписание на понедельник первой недели\n" +
			"> «/расписание 2 полное» - полное расписание второй недели"
		weekNum int
		dayName string
		err     error
	)

	params := strings.Split(body, " ")
	// Обязательно должно быть 2 параметра
	if len(params) != 2 {
		helpers.Logger.Warnw("Ошибка парсинга тела команды: количество параметров больше двух",
			"command", common.TimetableCommand,
			"body", body,
		)
		return 0, "", errors.New(wrongBody)
	}

	weekNum, err = strconv.Atoi(params[0])
	if err != nil {
		helpers.Logger.Warnw("Ошибка парсинга тела команды: первый параметр не число",
			"command", common.TimetableCommand,
			"body", body,
		)
		return 0, "", errors.New(wrongBody)
	}
	dayName = params[1]

	if weekNum != 1 && weekNum != 2 {
		helpers.Logger.Warnw("Ошибка парсинга тела команды: указан неверный номер недели",
			"command", common.TimetableCommand,
			"weekNum", weekNum,
		)
		return 0, "", errors.New("Ошибка в команде! Номер недели может принимать значения только 1 или 2")
	}
	if dayName == tools.Sunday {
		return 0, "", errors.New("Воскресенье - выходной!")
	}
	if !tools.IsStringInSlice(dayName, tools.Weekdays) && dayName != db.FullWeek {
		helpers.Logger.Warnw("Ошибка парсинга тела команды: указан неверный день недели",
			"command", common.TimetableCommand,
			"dayName", dayName,
		)
		return 0, "", errors.New("Указано неверное название дня недели")
	}

	return weekNum, dayName, nil
}

func TimetableCommand(d *Data) error {
	var (
		weekNum int    // номер недели
		dayName string // название дня недели
		err     error
	)
	now := tools.Now()

	t, err := db.GetTimetable()
	if err != nil {
		common.SendHandlerErrToSentry(d.command, err, common.DefaultHandlerBreadcrumbs(d.u, "", nil)...)
		return err
	}

	user, err := db.GetUserByVkID(d.u.Message.PeerID)
	if err != nil {
		common.SendHandlerErrToSentry(d.command, err, common.DefaultHandlerBreadcrumbs(d.u, "", nil)...)
		return err
	}

	switch body := strings.Trim(strings.ToLower(d.u.Message.MessageBody), " "); body {
	case "сегодня", "завтра":
		if body == "завтра" {
			now = now.Add(time.Hour * 24)
		}

		if int(now.Weekday()) == 0 {
			d.Answer = "В воскресенье нет занятий!"
			return nil
		}

		weekNum = tools.GetWeekNum(now)
		dayName = tools.GetWeekdayName(now)
	case "следующая":
		d.Answer, err = NextLesson(t, user.Group.Code)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	// Тело запроса вида: <weekNum> <dayName>. Прим: 1 понедельник, 2 пятница и тд.
	default:
		weekNum, dayName, err = parseBody(body)
		if err != nil {
			d.Answer = err.Error()
			return nil
		}
	}

	d.Answer, err = t.GetStringifyedSchedule(user.Group.Code, weekNum, dayName)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
