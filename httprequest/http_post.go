package httprequest

import (
	"bot_tg/envhandler"
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

func HttpPost(url string, data1 string, data2 string, headers map[string]string) (*http.Response, error) {
	var req *http.Request
	var err error

	// Проверяем необходимость создания multipart/form-data
	if headers["Content-Type"] == "multipart/form-data" {
		// Создаем временный буфер для multipart тела
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Добавляем первый файл (data1)
		if data1 != "" {
			filePart1, err := writer.CreateFormFile("file1", "file1.jpg") // "file1" - имя поля формы
			if err != nil {
				return nil, err
			}

			fileData1, err := ioutil.ReadFile(data1)
			if err != nil {
				return nil, err
			}
			filePart1.Write(fileData1)
		}

		// Добавляем второй файл (data2)
		if data2 != "" {
			filePart2, err := writer.CreateFormFile("file2", "file2.jpg") // "file2" - имя поля формы
			if err != nil {
				return nil, err
			}

			fileData2, err := ioutil.ReadFile(data2)
			if err != nil {
				return nil, err
			}
			filePart2.Write(fileData2)
		}

		// Закрываем writer для получения boundary
		writer.Close()

		// Создаем запрос с multipart телом
		req, err = http.NewRequest("POST", url, body)
		if err != nil {
			return nil, err
		}

		// Устанавливаем правильный Content-Type с boundary
		req.Header.Set("Content-Type", writer.FormDataContentType())
	} else {
		// Обычный запрос с текстовым телом
		req, err = http.NewRequest("POST", url, bytes.NewBufferString(data1)) // data1 используется как текстовое тело
		if err != nil {
			return nil, err
		}

		// Устанавливаем Content-Type из заголовков или по умолчанию
		if contentType, ok := headers["Content-Type"]; ok {
			req.Header.Set("Content-Type", contentType)
		} else {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	// Добавляем остальные заголовки
	for key, value := range headers {
		if key != "Content-Type" { // Content-Type уже установлен
			req.Header.Set(key, value)
		}
	}

	// Загружаем конфигурацию
	cfg, errEnv := envhandler.LoadConfig()
	if errEnv != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", errEnv)
	}

	// Создаем куки
	cookies := []*http.Cookie{
		{
			Name:  "user_id",
			Value: "7d40ed97-f964-4b49-a822-4656a57c3a6d", // cfg.UserId(), если есть такой метод
		},
		{
			Name:  "NUXT_LOCALE",
			Value: "ru",
		},
		{
			Name:  "TawkConnectionTime",
			Value: "0",
		},
		{
			Name:  "twk_uuid_66f0cb3de5982d6c7bb2f3cb",
			Value: cfg.TwlUuid, // Убедитесь, что это поле существует в конфигурации
		},
	}

	// Добавляем куки в запрос
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
