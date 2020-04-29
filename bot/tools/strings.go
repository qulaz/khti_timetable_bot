package tools

import (
	"strings"
)

// Удаление двух и более пробелов из текста. Поддерживаются многострочные тексты
func RemoveWhitespaces(s string) string {
	sliceOfStrings := strings.Split(s, "\n")
	clearSliceOfStrings := make([]string, 0, len(sliceOfStrings))

	for _, str := range sliceOfStrings {
		clearSliceOfStrings = append(clearSliceOfStrings, strings.Join(strings.Fields(str), " "))
	}

	return strings.Join(clearSliceOfStrings, "\n")
}

// Из двух переданных строк возвращает одну непустую строку. Если обе строки непустые, то приоритет отдается первой.
// В случае, если обе строки пустые, возвращается пустая строка
func SelectNonEmptyString(s1, s2 string) string {
	if len(strings.TrimSpace(s1)) != 0 {
		return s1
	}
	if len(strings.TrimSpace(s2)) != 0 {
		return s2
	}
	return ""
}

// Возвращает правильную форму существительного с числительным.
//  number - числительное для которого нужно подобрать правильную форму слова
//  titles - слайс со словом в 3 формах в правильной последовательности: именительный падеж единственное число
//  (сочетается с 1), именительный падеж множественное число (сочетается с 2),
//  родительный падеж множественное число (сочетается с 5)
// Пример:
//  titles := []string{"минута", "минуты", "минут"}
//  SelectionNounForm(1, titles) // минута
//  SelectionNounForm(2, titles) // минуты
//  SelectionNounForm(5, titles) // минут
//
// сс: https://gist.github.com/chiliec/22e34af2d08a964fc1418908a19b0c15
func SelectionNounForm(number int, titles []string) string {
	cases := []int{2, 0, 1, 1, 1, 2}
	var currentCase int
	if number%100 > 4 && number%100 < 20 {
		currentCase = 2
	} else if number%10 < 5 {
		currentCase = cases[number%10]
	} else {
		currentCase = cases[5]
	}
	return titles[currentCase]
}
