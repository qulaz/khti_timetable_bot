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
