package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Определяем структуру проекта
	projectStructure := map[string][]string{
		"pgmonitor/config":      {"config.yaml", "config.go"},
		"pgmonitor/internal/collector": {"database.go", "queries.go", "locks.go", "resources.go", "execution.go"},
		"pgmonitor/internal/models":     {"models.go"},
		"pgmonitor/internal/repository": {"repository.go"},
		"pgmonitor/pkg/logger":         {"logger.go"},
		"pgmonitor/pkg/utils":          {"utils.go"},
		"pgmonitor":                    {"go.mod", "go.sum", "main.go"},
	}

	// Создаем директории и файлы
	for dir, files := range projectStructure {
		// Создаем директорию
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Ошибка при создании директории %s: %v\n", dir, err)
			continue
		}

		// Создаем файлы
		for _, file := range files {
			filePath := filepath.Join(dir, file)
			if _, err := os.Create(filePath); err != nil {
				fmt.Printf("Ошибка при создании файла %s: %v\n", filePath, err)
			}
		}
	}

	fmt.Println("Структура проекта успешно создана!")
}
