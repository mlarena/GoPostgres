package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"pgmonitor/config"
	"pgmonitor/internal/collector"
	"pgmonitor/internal/repository"
	"pgmonitor/pkg/logger"
	
	_ "github.com/lib/pq"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	// Инициализация логгера
	appLogger, err := logger.NewLogger(&cfg.Logging) // Теперь типы совместимы
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer appLogger.Close()
	
	appLogger.Info("Starting PostgreSQL monitor")
	
	// Подключение к базам данных
	workDB, err := sql.Open("postgres", cfg.Database["workdb"].ConnectionString())
	if err != nil {
		appLogger.Fatalf("Failed to connect to workdb: %v", err)
	}
	defer workDB.Close()
	
	sourceDB, err := sql.Open("postgres", cfg.Database["source"].ConnectionString())
	if err != nil {
		appLogger.Fatalf("Failed to connect to source db: %v", err)
	}
	defer sourceDB.Close()
	
	// Инициализация репозиториев и коллекторов
	repo := repository.NewRepository(workDB, appLogger)
	monitor := collector.NewCollector(repo, appLogger, cfg)
	
	// Обработка сигналов для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// Запуск мониторинга
	go func() {
		checkInterval, err := cfg.Monitoring.GetCheckInterval()
		if err != nil {
			appLogger.Errorf("Invalid check interval: %v", err)
			return
		}
		
		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				if err := monitor.Run(ctx); err != nil {
					appLogger.Errorf("Monitoring failed: %v", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	
	// Ожидание сигнала завершения
	<-sigChan
	appLogger.Info("Shutting down monitor...")
	cancel()
	time.Sleep(1 * time.Second) // Даем время для завершения операций
	appLogger.Info("Monitor stopped")
}