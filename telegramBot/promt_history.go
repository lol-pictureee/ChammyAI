package telegramBot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func WriteMessageToFile(update tgbotapi.Update, userDir string) error {

	filename := filepath.Join(userDir, "messages.txt")

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	now := time.Now()
	timestamp := now.Format("02.01.2006, 15:04:05")
	userID := strconv.FormatInt(update.Message.From.ID, 10)

	messageText := fmt.Sprintf("[User: %s - (id: %s)] [Date: %s], [Prompt: %s]\n", update.Message.From.FirstName, userID, timestamp, update.Message.Text)

	_, err = file.WriteString(messageText)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
