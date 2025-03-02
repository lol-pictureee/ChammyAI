package bot_tg

import (
	"log/slog"
	"os"
	"testing"
)

func TestName(t *testing.T) {
	// Инициализация slog с форматом JSON и указанием места возникновения ошибки.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug, // Уровень логирования (Debug, Info, Warn, Error)
		AddSource: true,            // Добавлять информацию об источнике лога (файл, строка)
	}))

	// Установка логгера по умолчанию.  Теперь slog.Info, slog.Error и т.д. будут использовать этот logger.
	slog.SetDefault(logger)
	slog.Warn(
		"Ошибка генерации изображения",
		"user", "Chamik",
		"Error", "Not Found",
	)

}
