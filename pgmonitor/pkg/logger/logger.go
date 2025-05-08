package logger

import (
	"fmt"
	"os"
	"path/filepath"
	
	"pgmonitor/config"
	
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	*zap.SugaredLogger
}

// Изменяем сигнатуру функции для использования config.LoggingConfig
func NewLogger(cfg *config.LoggingConfig) (*Logger, error) {
	// Создаем директорию для логов если не существует
	if err := os.MkdirAll(cfg.Path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}
	
	logFile := filepath.Join(cfg.Path, "pgmonitor.log")
	
	// Настройка уровней логирования
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}
	
	// Настройка writer для ротации логов
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   true,
	})
	
	// Настройка encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	
	// Ядро логирования
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writer, zapcore.AddSync(os.Stdout)),
		level,
	)
	
	// Создаем логгер
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	sugar := logger.Sugar()
	
	return &Logger{sugar}, nil
}

func (l *Logger) Close() error {
	return l.Sync()
}