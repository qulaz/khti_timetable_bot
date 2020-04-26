package helpers

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func LoadBytesFromFile(file string) []byte {
	path := filepath.Join("testdata", file)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		dir, _ := os.Getwd()
		log.Fatalf("Не удалось загрузить тестовые данные %s: %s", dir, err)
	}
	return bytes
}

func LoadStringFromFile(file string) string {
	return string(LoadBytesFromFile(file))
}
