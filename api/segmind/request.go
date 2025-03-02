package segmind

import (
	"bot_tg/converter"
	"bot_tg/envhandler"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func RequestSwapService(SourceUrl string, TargetUrl string) error {
	cfg, errEnv := envhandler.LoadConfig()
	if errEnv != nil {
		fmt.Errorf("ошибка загрузки конфигурации: %w", errEnv)

	}

	urlApi := "https://api.segmind.com/v1/faceswap-v2" //////////////////////////////////////////////
	headers := Headers{
		APIKey: cfg.ApiKeySegmind,
	}
	data := Data{

		SourceImg:          converter.ImageUrlToBase64(SourceUrl),
		TargetImg:          converter.ImageUrlToBase64(TargetUrl),
		InputFacesIndex:    0,
		SourceFacesIndex:   0,
		FaceRestore:        "codeformer-v0.1.0.pth",
		Interpolation:      "Bilinear",
		DetectionFaceOrder: "large-small",
		FaceDetection:      "retinaface_resnet50",
		DetectGenderInput:  "no",
		DetectGenderSource: "no",
		FaceRestoreWeight:  0.75,
		ImageFormat:        "jpeg",
		ImageQuality:       95,
		Base64:             false,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных: %v", err)
	}

	req, err := http.NewRequest("POST", urlApi, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", headers.APIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Код ответа:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ошибка, ответ от сервера: %s", resp.Status)
	}

	fmt.Println("Запрос выполнен успешно.")

	if resp.StatusCode == http.StatusOK {
		outFile, err := os.Create("output_image.jpeg")
		if err != nil {
			return fmt.Errorf("ошибка создания файла: %v", err)
		}
		defer outFile.Close()
		_, err = io.Copy(outFile, resp.Body)
		if err != nil {
			return fmt.Errorf("ошибка записи в файл: %v", err)
		}

		fmt.Println("Изображение успешно сохранено как output_image.jpeg.")
	} else {
		return fmt.Errorf("неуспешный ответ от сервера: %s", resp.Status)
	}
	return nil
}
