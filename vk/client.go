package vk

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

const DEFAULT_TIMEOUT = 30

// Создание http.Client со стандартными параметрами
func NewClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * DEFAULT_TIMEOUT,
	}
}

// Get-запрос
func getRequest(client *http.Client, url string) ([]byte, error) {
	response, err := client.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, errors.Wrap(err, "Ошибка чтения тела ответа")
	}

	return body, err
}
