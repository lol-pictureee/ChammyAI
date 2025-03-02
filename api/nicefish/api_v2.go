package nicefish

import (
	"bot_tg/envhandler"
	"bot_tg/routes"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func DownloadPhoto(requestID string) (string, error) {
	cfg, errEnv := envhandler.LoadConfig()
	if errEnv != nil {
		return "", fmt.Errorf("ошибка загрузки конфигурации: %w", errEnv)
	}

	headers := map[string]string{
		"authority":       "aifaceswapper.io",
		"method":          "GET",
		"path":            fmt.Sprintf("/api/nicefish/fs/result?request_id=%v", requestID),
		"scheme":          "https",
		"accept":          "*/*",
		"accept-encoding": "gzip, deflate, br, zstd",
		"accept-language": "ru-RU,ru;q=0.5",
		"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
		"authorization":   "e3e87194-99c0-4c63-b17b-146ee70e0d0b",
	}

	cookies := Cookies(cfg.TwlUuid)

	urlResult := "https://aifaceswapper.io/api/nicefish/fs/result?request_id=" + requestID

	req, err := http.NewRequest("GET", urlResult, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return "", fmt.Errorf("ошибка создания gzip reader: %w", err)
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", fmt.Errorf("ошибка декодирования JSON: %w. Body: %s", err, string(body))
	}

	dataMap, ok := data["data"].(map[string]interface{})
	if !ok {
		log.Println("Не удалось получить доступ к data")
		log.Println("Ответ сервера:", string(body))
		return "", fmt.Errorf("не удалось получить доступ к data")
	}

	photoURL, ok := dataMap["result_img_url"].(string)
	if !ok {
		log.Println("Не удалось извлечь result_img_url из data")
		log.Println("Ответ сервера:", string(body))
		return "", fmt.Errorf("не удалось извлечь result_img_url из data")
	}

	return photoURL, nil

	return "", fmt.Errorf("не удалось извлечь URL фото из ответа")
}

func RequestID(targetImage string, targetFace string, userID string) (string, error) {

	cfg, errEnv := envhandler.LoadConfig()
	if errEnv != nil {
		fmt.Errorf("ошибка загрузки конфигурации: %w", errEnv)

	}

	headers := map[string]string{
		"authority":        "aifaceswapper.io",
		"method":           "POST",
		"path":             "/api/nicefish/fs/singleface",
		"scheme":           "https",
		"accept":           "*/*",
		"accept-encoding":  "gzip, deflate, br, zstd",
		"accept-language":  "ru-RU,ru;q=0.5",
		"authorization":    "e3e87194-99c0-4c63-b17b-146ee70e0d0b",
		"origin":           "https://aifaceswapper.io",
		"referer":          "https://aifaceswapper.io/ru",
		"sec-ch-ua-mobile": "?0",
		"sec-fetch-dest":   "empty",
		"sec-fetch-mode":   "cors",
		"sec-fetch-site":   "same-origin",
		"sec-gpc":          "1",
		"user-agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
	}

	cookies := Cookies(cfg.TwlUuid)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	files := map[string]string{
		"target_image_file": targetImage,
		"target_face_file":  targetFace,
	}

	for fieldName, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			return "", fmt.Errorf("ошибка открытия файла %s: %w", filePath, err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile(fieldName, filePath)
		if err != nil {
			return "", fmt.Errorf("ошибка создания части формы для файла %s: %w", filePath, err)
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return "", fmt.Errorf("ошибка копирования данных файла %s: %w", filePath, err)
		}
	}

	writer.Close()

	req, err := http.NewRequest("POST", routes.ApiFsNiceFish.String(), body)
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.StatusCode)
	contentEncoding := resp.Header.Get("Content-Encoding")

	var reader io.ReadCloser
	reader = resp.Body

	if contentEncoding == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return "", fmt.Errorf("ошибка при создании gzip reader. Проверьте contentEncoding: %w", err)
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	respBody, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(respBody, &data)
	if err != nil {

		return "", err
	}

	//fmt.Println("response body:", data)

	if code, ok := data["code"].(float64); ok && code == 100000 {
		data := data["data"].(map[string]interface{})
		requestID := data["request_id"].(string)
		return requestID, nil
		//fmt.Println("request id:", requestID)

	}
	return "", fmt.Errorf("не удалось получить request_id из ответа")

}
