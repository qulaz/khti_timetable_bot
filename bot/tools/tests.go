package tools

import (
	"io/ioutil"
	"log"
	"testing"
)

func LoadBytesFromFile(path string) []byte {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Не удалось загрузить тестовые данные: %s", err)
	}
	return bytes
}

func LoadStringFromFile(path string) string {
	return string(LoadBytesFromFile(path))
}

// Функция, позволяющая зафейлить тест в случае, если assert функция вернет false
//  import (
//      "testing"
//      "fmt"
//      "github.com/stretchr/testify/assert"
//  )
//
//  err := fmt.Errorf("test %s", "error")
//  Fatal(t, assert.NoError(t, err))
func Fatal(t *testing.T, assertion bool) {
	if !assertion {
		t.Fatal()
	}
}
