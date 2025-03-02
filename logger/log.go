package logger

import (
	"io"
	"log/slog"
	"os"
	"time"
)

func LogS() {

	logFile, err := os.OpenFile("app.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		slog.Error("Не удалось открыть файл логов", "error", err)
		os.Exit(1)
	}
	defer logFile.Close()

	// Создаем MultiWriter для записи одновременно в файл и в консоль.
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Инициализация slog с форматом JSON, форматированным временем, записью в файл и консоль.
	logger := slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format(time.DateTime))
			}
			return a
		},
	}))
	slog.SetDefault(logger)

}
