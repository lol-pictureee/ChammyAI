package envhandler

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Cookie             string
	BotToken           string
	ApiKeyUrlUpload    string
	ApiKeyFuseBrain    string
	ApiSecretFuseBrain string
	ApiKeyServiceThree string
	ApiKeySegmind      string
	TwlUuid            string
	// soon...
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("config.env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cookie := os.Getenv("Cookie")
	if cookie == "" {
		return nil, fmt.Errorf("cookie is not set")
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("BotToken is not set")
	}

	apiKeyUrlUpload := os.Getenv("API_KEY_URL_UPLOAD")
	if apiKeyUrlUpload == "" {
		return nil, fmt.Errorf("ApiKeyUrlUpload is not set")
	}

	apiKeyFuseBrain := os.Getenv("API_KEY_FUSE_BRAIN")
	if apiKeyFuseBrain == "" {
		return nil, fmt.Errorf("ApiKeyFuseBrain is not set")
	}

	apiSecretFuseBrain := os.Getenv("API_SECRET_FUSE_BRAIN")
	if apiSecretFuseBrain == "" {
		return nil, fmt.Errorf("API_SECRET_FUSE_BRAIN is not set")
	}

	apiKeySegmind := os.Getenv("API_KEY_SEGMIND")
	if apiKeySegmind == "" {
		return nil, fmt.Errorf("API_KEY_SEGMIND is not set")
	}

	twlUid := os.Getenv("TWL_UUID")
	if twlUid == "" {
		return nil, fmt.Errorf("TWL_UUID is not set")
	}

	return &Config{
		Cookie:             cookie,
		BotToken:           botToken,
		ApiKeyUrlUpload:    apiKeyUrlUpload,
		ApiKeyFuseBrain:    apiKeyFuseBrain,
		ApiSecretFuseBrain: apiSecretFuseBrain,
		ApiKeySegmind:      apiKeySegmind,
		TwlUuid:            twlUid,
	}, nil
}
