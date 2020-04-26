package vk

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMessage_SendOptionsValidate(t *testing.T) {
	opt_without_reciever := &SendOptions{Message: "Hello"}
	if err := opt_without_reciever.validate(); err == nil {
		t.Error("Пропущены параметры отправки с не заданным получателем")
	} else {
		if !strings.Contains(fmt.Sprint(err), "не задан получатель сообщения") {
			t.Error("1: Неверное сообщение об ошибке")
		}
	}

	opt_without_content := &SendOptions{PeerID: 123}
	if err := opt_without_content.validate(); err == nil {
		t.Error("Пропущен контент сообщения")
	} else {
		if !strings.Contains(fmt.Sprint(err), "не задан контент сообщения") {
			t.Error("2: Неверное сообщение об ошибке")
		}
	}

	opt_with_content_tb_processed := &SendOptions{
		PeerID:          123,
		Message:         "smth",
		Attachment:      []string{"photo123_123", "audio123_123", "video123_123"},
		ForwardMessages: []int{123, 12, 1},
	}
	opt_with_content_tb_processed.validate()

	assert.Equal(t, "photo123_123,audio123_123,video123_123", opt_with_content_tb_processed.AttachmentProcessed)
	assert.Equal(t, "123,12,1", opt_with_content_tb_processed.ForwardMessagesProcessed)
}

func TestMessage_SendOptionsValidateSendMultiple(t *testing.T) {
	// Флаг IsSendMultiple позволяет указывать, опции это для одиночного сообщения или для рассылки на несколько
	// пользователей. Такое разделение происходит из-за того, что ВК при отправке одиночного сообщения и рассылки
	// использует совершенно разные ответы. Чтобы была нормальная типизация, было принято решение сдлать отдельную
	// функцию для отправки одиночных сообщений и массовой рассылки. В случае одиночного сообщения, отправитель
	// указывается в SendOptions, а в случае с массовой рассылкой, ID пользователей указываются в аргументах ф-ии
	opt_without_reciever_but_with_multiple := &SendOptions{
		Message:        "Hello",
		IsSendMultiple: true,
	}
	if err := opt_without_reciever_but_with_multiple.validate(); err != nil {
		t.Error("Установлен флаг IsSendMultiple, ошибки валидации быть не должно!")
	}
}

func TestMessage_SendOptionsToParams(t *testing.T) {
	expectedValue := H{
		"keyboard": "{\"one_time\":false,\"buttons\":[" +
			"[{\"action\":{\"type\":\"text\",\"label\":\"\"},\"color\":\"\"}]," +
			"[{\"action\":{\"type\":\"text\",\"label\":\"\"},\"color\":\"\"}]," +
			"[{\"action\":{\"type\":\"text\",\"label\":\"\"},\"color\":\"\"}]" +
			"],\"inline\":false}",
		"message":          "message",
		"lat":              "15.000000",
		"long":             "16.000000",
		"attachment":       "photo123_123,audio123_123,video123_123",
		"forward_messages": "123,12,1",
		"peer_id":          "123",
	}
	k := NewKeyboard(false)
	k.AddTextButton("", "", nil)
	k.AddRow()
	k.AddTextButton("", "", nil)
	k.AddRow()
	k.AddTextButton("", "", nil)

	opt := &SendOptions{
		PeerID:          123,
		Message:         "message",
		Lat:             CreateFloat64(15),
		Long:            CreateFloat64(16),
		Attachment:      []string{"photo123_123", "audio123_123", "video123_123"},
		ForwardMessages: []int{123, 12, 1},
		Keyboard:        k,
	}
	opt.validate()

	result, err := opt.toParams()
	assert.NoError(t, err)

	expectedValue["random_id"] = result["random_id"]

	bexpectedValue, err := json.Marshal(expectedValue)
	assert.NoError(t, err, "Ошибка сериализации")

	bresult, err := json.Marshal(result)
	assert.NoError(t, err, "Ошибка сериализации")

	assert.Equal(t, string(bexpectedValue), string(bresult))
}
