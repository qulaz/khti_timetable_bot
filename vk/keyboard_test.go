package vk

import (
	"strconv"
	"testing"
)

func TestKeyboard_AddRow(t *testing.T) {
	// В обычной клавиатуре может быть максимум 10 строк
	k := NewKeyboard(false)

	for i := 0; i < 10; i++ {
		if i == 0 {
			if err := k.AddLinkButton("", ""); err != nil {
				t.Errorf("Ошибка при добавлении первой кнопки в обычную клавиатуру: %+v\n", err)
			}
		} else {
			if err := k.AddLinkButton("", ""); err != nil {
				t.Errorf("Ошибка при добавлении %d кнопки в обычную клавиатуру: %+v\n", i+1, err)
			}
			if err := k.AddRow(); err != nil {
				t.Errorf("Ошибка при добавлении %d строки в обычную клавиатуру: %+v\n", i+1, err)
			}
		}
	}
	if err := k.AddRow(); err == nil {
		t.Error("Нет ошибки при добавлении 11 строки в обычную клавиатуру")
	}
	if len(k.Buttons) > 10 {
		t.Error("В обычной клавиатуре больше 10 строк")
	}

	// В inline клавиатуре может быть максимум 6 строк
	k = NewInlineKeyboard()

	for i := 0; i < 6; i++ {
		if i == 0 {
			if err := k.AddLinkButton("", ""); err != nil {
				t.Errorf("Ошибка при добавлении первой кнопки в inline клавиатуру: %+v\n", err)
			}
		} else {
			if err := k.AddLinkButton("", ""); err != nil {
				t.Errorf("Ошибка при добавлении %d кнопки в inline клавиатуру: %+v\n", i+1, err)
			}
			if err := k.AddRow(); err != nil {
				t.Errorf("Ошибка при добавлении %d строки в inline клавиатуру: %+v\n", i+1, err)
			}
		}
	}
	if err := k.AddRow(); err == nil {
		t.Error("Нет ошибки при добавлении 6 строки в inline клавиатуру")
	}
	if len(k.Buttons) > 10 {
		t.Error("В обычной inline больше 6 строк")
	}
}

func TestKeyboard_AddButton_LimitButtonsCount(t *testing.T) {
	// Тестирование на превышение кол-ва кнопок в клавиатуре
	k := NewKeyboard(false)

	for i := 0; i < 40; i++ {
		if err := k.AddTextButton(strconv.Itoa(i), COLOR_PRIMARY, nil); err != nil {
			t.Errorf("Ошибка при добавлении %d кнопки в обычную клавиатуру: %+v\n", i+1, err)
		}
	}
	if err := k.AddTextButton(strconv.Itoa(40), COLOR_PRIMARY, nil); err == nil {
		t.Error("Нет ошибки при добавлении 40 кнопки в обычную клавиатуру")
	}
	if k.ButtonCount() > 40 {
		t.Error("В обычной клавиатуре больше 40 кнопок")
	}

	k = NewInlineKeyboard()

	for i := 0; i < 10; i++ {
		if err := k.AddTextButton(strconv.Itoa(i), COLOR_PRIMARY, nil); err != nil {
			t.Errorf("Ошибка при добавлении %d кнопки в inline клавиатуру: %+v\n", i+1, err)
		}
	}
	if err := k.AddTextButton(strconv.Itoa(10), COLOR_PRIMARY, nil); err == nil {
		t.Error("Нет ошибки при добавлении 10 кнопки в inline клавиатуру")
	}
	if k.ButtonCount() > 10 {
		t.Error("В inline клавиатуре больше 10 кнопок")
	}
}

func TestKeyboard_AddButton_AutoAddRowForSpecialButtons(t *testing.T) {
	// В одной строке может быть до 5 кнопок, но есть специальные кнопки, которые должны занимать всю ширину строки, те
	// быть одними в строке
	k := NewKeyboard(false)

	k.AddTextButton("sample text", COLOR_PRIMARY, nil)
	k.AddTextButton("sample text", COLOR_PRIMARY, nil)
	k.AddTextButton("sample text", COLOR_PRIMARY, nil)
	k.AddTextButton("sample text", COLOR_PRIMARY, nil)

	k.AddVkPayButton("")

	k.AddTextButton("sample text", COLOR_PRIMARY, nil)
	k.AddTextButton("sample text", COLOR_PRIMARY, nil)
	k.AddTextButton("sample text", COLOR_PRIMARY, nil)
	k.AddTextButton("sample text", COLOR_PRIMARY, nil)

	k.AddLocationButton(nil)

	k.AddLinkButton("", "")
	k.AddLinkButton("", "")
	k.AddRow()

	k.AddTextButton("sample text", COLOR_PRIMARY, nil)
	k.AddTextButton("sample text", COLOR_PRIMARY, nil)

	k.AddVkPayButton("")

	k.AddVkPayButton("")

	k.AddTextButton("sample text", COLOR_PRIMARY, nil)
	k.AddTextButton("sample text", COLOR_PRIMARY, nil)
	k.AddTextButton("sample text", COLOR_PRIMARY, nil)
	k.AddTextButton("sample text", COLOR_PRIMARY, nil)
	k.AddTextButton("sample text", COLOR_PRIMARY, nil)

	k.AddVkPayButton("")

	if err := k.AddTextButton("sample text", COLOR_PRIMARY, nil); err == nil {
		t.Error("Добавилась кнопка, которой быть не должно", len(k.Buttons))
	}
}

func TestKeyboard_AddRow_BlankRows(t *testing.T) {
	k := NewKeyboard(true)
	k.AddTextButton("", "", nil)
	k.AddLocationButton(nil)
	k.AddRow()
	k.AddLinkButton("", "")
	k.AddVkPayButton("")
	k.AddRow()
	k.AddTextButton("test3", COLOR_PRIMARY, nil)
	k.AddVkAppsButton(0, 0, "", "")

	for i, row := range k.Buttons {
		// последняя пустая строка обрабатывается ф-ей клавиатуры toJson()
		if len(row) == 0 && (len(k.Buttons)-1) != i {
			t.Error("1 Обнаружена пустая строка ", i+1)
		}
	}
}

func TestKeyboard_AddButton_AutoAddRowForSpecialButtonsInline(t *testing.T) {
	// В одной строке может быть до 5 кнопок, но есть специальные кнопки, которые должны занимать всю ширину строки, те
	// быть одними в строке
	k := NewInlineKeyboard()

	k.AddTextButton("sample text", COLOR_PRIMARY, nil)
	k.AddTextButton("sample text", COLOR_PRIMARY, nil)
	k.AddTextButton("sample text", COLOR_PRIMARY, nil)
	k.AddTextButton("sample text", COLOR_PRIMARY, nil)

	k.AddVkPayButton("")

	k.AddTextButton("sample text", COLOR_PRIMARY, nil)
	k.AddTextButton("sample text", COLOR_PRIMARY, nil)
	k.AddTextButton("sample text", COLOR_PRIMARY, nil)
	k.AddTextButton("sample text", COLOR_PRIMARY, nil)

	k.AddLocationButton(nil)

	k.AddLinkButton("", "")
	k.AddLinkButton("", "")
	k.AddRow()

	k.AddVkPayButton("")

	if err := k.AddTextButton("sample text", COLOR_PRIMARY, nil); err == nil {
		t.Error("Добавилась кнопка, которой быть не должно", len(k.Buttons))
	}
}

func TestKeyboard_AddButton_Test1(t *testing.T) {
	k := NewKeyboard(true)
	for i := 0; i < 27; i++ {
		if i%3 == 0 {
			if err := k.AddRow(); err != nil {
				t.Error("не добавилась строка")
			}
		}
		if err := k.AddTextButton("a", COLOR_PRIMARY, &ButtonPayload{}); err != nil {
			t.Error("не добавилась кнопка")
		}
	}
	if err := k.AddRow(); err != nil {
		t.Error("не добавилась строка после цикла")
	}
	if err := k.AddTextButton("a", COLOR_PRIMARY, &ButtonPayload{}); err != nil {
		t.Error("не добавилась кнопка после цикла")
	}
}
