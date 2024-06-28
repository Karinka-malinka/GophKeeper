// Пакет logger предоставляет простую настройку логирования с использованием пакета slog.
package logger

import (
	"log/slog"
	"os"
)

// SetUpLogger настраивает логгер по умолчанию с форматированием в JSON и уровнем отладки.
func SetUpLogger() {

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}

	// Создаем новый логгер с обработчиком JSON, записывающим в os.Stdout
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	// Устанавливаем созданный логгер как логгер по умолчанию
	slog.SetDefault(logger)
}
