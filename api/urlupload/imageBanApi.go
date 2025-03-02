package urlupload

import (
	"bot_tg/envhandler"
	"bot_tg/routes"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type Headers struct {
	APIKey string `json:"Authorization"`
}

type ImageBanResponse struct {
	Data struct {
		Link string `json:"link"`
	} `json:"data"`
}

func RequestUrlUpload(imagePath string) (string, error) {
	cfg, err := envhandler.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("Ошибка открытия файла: %w", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
	if err != nil {
		return "", fmt.Errorf("Ошибка создания формы файла: %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("Ошибка копирования файла", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("error closing multipart writer: %w", err)
	}

	req, err := http.NewRequest("POST", routes.ApiUrlUpload.String(), body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", cfg.ApiKeyUrlUpload)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(respBody))
	}

	var imageBanResponse ImageBanResponse
	err = json.Unmarshal(respBody, &imageBanResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling JSON response: %w", err)
	}

	return imageBanResponse.Data.Link, nil
}

func RunUrlUpload(imagePath string) string {
	imageUrl, err := RequestUrlUpload(imagePath)
	if err != nil {
		fmt.Println("Error:", err)

	}

	fmt.Println("Image URL:", imageUrl)
	return imageUrl
}
