package fusebrain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"
)

type Text2ImageAPI struct {
	URL       string
	APIKey    string
	APISecret string
	Client    *http.Client
}

type GenerateResponse struct {
	UUID   string `json:"uuid"`
	Status string `json:"status"`
}

type StatusResponse struct {
	Status string `json:"status"`
	Error  string `json:"errorDescription"`
	Result struct {
		Files    []string `json:"files"`
		Censored bool     `json:"censored"`
	} `json:"result"`
}

func NewText2ImageAPI(url, apiKey, apiSecret string) *Text2ImageAPI {
	return &Text2ImageAPI{
		URL:       url,
		APIKey:    apiKey,
		APISecret: apiSecret,
		Client:    &http.Client{},
	}
}

func (api *Text2ImageAPI) GetModel() (string, error) {
	req, err := http.NewRequest("GET", api.URL+"key/api/v1/pipelines", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-Key", "Key "+api.APIKey)
	req.Header.Set("X-Secret", "Secret "+api.APISecret)

	resp, err := api.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var models []struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		return "", err
	}

	if len(models) == 0 {
		return "", fmt.Errorf("no models available")
	}

	return models[0].ID, nil
}

func (api *Text2ImageAPI) Generate(prompt string, model string, images, width, height int) (string, error) {
	params := map[string]interface{}{
		"type":      "GENERATE",
		"numImages": images,
		"width":     width,
		"height":    height,
		"generateParams": map[string]string{
			"query": prompt,
		},
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	partHeader := textproto.MIMEHeader{}
	partHeader.Set("Content-Disposition", `form-data; name="params"`)
	partHeader.Set("Content-Type", "application/json")
	paramsPart, err := writer.CreatePart(partHeader)
	if err != nil {
		return "", err
	}
	paramsPart.Write(paramsJSON)

	writer.WriteField("pipeline_id", fmt.Sprintf("%s", model))
	writer.Close()

	req, err := http.NewRequest("POST", api.URL+"key/api/v1/pipeline/run", body)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-Key", "Key "+api.APIKey)
	req.Header.Set("X-Secret", "Secret "+api.APISecret)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := api.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Status != "INITIAL" {
		return "", fmt.Errorf("generation failed: %s", result.Status)
	}

	return result.UUID, nil
}

func (api *Text2ImageAPI) CheckGeneration(requestID string, attempts int, delay time.Duration) ([]string, error) {
	url := api.URL + "key/api/v1/pipeline/status/" + requestID

	for i := 0; i < attempts; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("X-Key", "Key "+api.APIKey)
		req.Header.Set("X-Secret", "Secret "+api.APISecret)

		resp, err := api.Client.Do(req)
		if err != nil {
			return nil, err
		}

		var status StatusResponse
		if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		switch status.Status {
		case "DONE":
			return status.Result.Files, nil
		case "FAIL":
			return nil, fmt.Errorf("generation failed: %s (error: %s)", status.Error)
		}

		time.Sleep(delay)
	}

	return nil, fmt.Errorf("maximum attempts reached")
}
