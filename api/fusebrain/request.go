package fusebrain

import (
	"bot_tg/converter"
	"bot_tg/envhandler"
	"fmt"
	"github.com/google/uuid"
	"path/filepath"
	"time"
)

func Request(prompt string, userID string) (string, error) {
	cfg, err := envhandler.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки конфигурации: %w", err) // Возвращаем ошибку, не завершаем программу
	}

	api := NewText2ImageAPI(
		"https://api-key.fusionbrain.ai/",
		cfg.ApiKeyFuseBrain,
		cfg.ApiSecretFuseBrain,
	)
	modelID, err := api.GetModel()
	if err != nil {
		return "", fmt.Errorf("ошибка получения modelID: %w", err) // Возвращаем ошибку
	}
	fullPrompt := "Стиль: Всегда человек с косоглазием, реализм. Лицо смотрит прямо в кадр на зрителя, " + prompt
	uuidApi, err := api.Generate(fullPrompt, modelID, 1, 1024, 1024)
	if err != nil {
		return "", fmt.Errorf("ошибка генерации изображения: %w", err) // Возвращаем ошибку
	}

	images, err := api.CheckGeneration(uuidApi, 25, 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("ошибка проверки генерации: %w", err) // Возвращаем ошибку
	}

	if len(images) > 0 {

		userDir := filepath.Join("users", userID)

		uuidGen := uuid.NewString()
		imageName := fmt.Sprintf("%s.jpg", uuidGen)
		outputPath := filepath.Join(userDir, imageName)

		err := converter.DecodeBase64Image(images[0], outputPath)

		if err != nil {
			return "", fmt.Errorf("ошибка декодирования base64 изображения: %w", err) // Возвращаем ошибку
		}

		fmt.Println("Изображение успешно декодировано и сохранено в:", outputPath)
		return outputPath, nil

	} else {
		return "", fmt.Errorf("не было получено ни одного изображения")
	}
}
