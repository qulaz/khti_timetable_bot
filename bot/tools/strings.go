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
