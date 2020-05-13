package parser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
	"gitlab.com/qulaz/khti_timetable_bot/bot/service"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var baseUrl = "http://khti.sfu-kras.ru"

// Возвращает ссылку на расписание с сайта ХТИ
func getTimetableLink() (string, error) {
	var (
		link string
		err  = errors.New("Ссылка на расписание не найдена")
	)

	res, err := http.Get(baseUrl + "/obuchenie/raspisanie-zanyatiy.php")
	if err != nil {
		return "", errors.Wrap(err, "Ошибка получения страницы с расписанием")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		helpers.Logger.Warnw(
			"Ошибка получения страницы с расписанием",
			"statusCode", res.StatusCode, "body", res.Body, "res", res,
		)
		return "", errors.Errorf("Ошибка получения страницы с расписанием: statusCode: %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "Ошибка парсинга страницы с расписанием")
	}

	doc.Find(".s-c-content p").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if text := s.Text(); strings.Contains(text, "Расписание занятий") {
			if relativeLink, ok := s.Find("a").Attr("href"); ok {
				link = baseUrl + relativeLink
				err = nil
			}
			return false
		}
		return true
	})

	return link, err
}

// Скачивание файла
func DownloadFile(filePath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return errors.WithStack(err)
}

// Скачивание актуального расписания с сайта ХТИ. Возвращает путь до скачанного файла расписания
func GrabTimetableFromSite() (string, error) {
	url, err := getTimetableLink()
	if err != nil {
		return "", errors.WithStack(err)
	}

	fileName := "timetable" + filepath.Ext(url)
	err = DownloadFile(fileName, url)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return fileName, nil
}

// Функция, которая переодически выполняется и проверяет изменилось ли расписание на сайте. В случае, если нашлось
// новое расписание - оно записывается в базу данных и всем пользователям отправляется уведомление
func UpdateTimetable(b *vk.Bot) error {
	helpers.Logger.Info("Запустилась задача проверки обновления расписания")

	filePath, err := GrabTimetableFromSite()
	if err != nil {
		return errors.Wrap(err, "Ошибка получения расписания с сайта")
	}

	t, err := Parse(filePath)
	if err != nil {
		return errors.Wrap(err, "Ошибка парсинга файла расписания")
	}

	isTimetableExists, err := db.IsTimetableExists()
	if err != nil {
		return errors.WithStack(err)
	}
	if !isTimetableExists {
		helpers.Logger.Info("В БД еще нет расписания. Записываем полученное")
		if err := t.WriteInDB(); err != nil {
			return errors.Wrap(err, "Ошибка записи новго расписания в БД")
		}
		return nil
	}

	isNewTimetable, err := t.IsNewTimetable()
	if err != nil {
		return errors.Wrap(err, "Ошибка сравнения расписаний")
	}
	if isNewTimetable {
		helpers.Logger.Info("Найденное новое расписание. Записываем его в БД")
		if err := t.WriteInDB(); err != nil {
			return errors.Wrap(err, "Ошибка записи новго расписания в БД")
		}
		if err := service.SendNotifyAboutTimetableUpdate(b); err != nil {
			return errors.Wrap(err, "Ошибка отправки уведомлений")
		}
	} else {
		helpers.Logger.Info("Расписание не обновилось")
	}

	helpers.Logger.Info("Успешно завершилась задача проверки обновления расписания")
	return nil
}
