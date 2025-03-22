package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Создаем временный конфиг для теста
	tmpConfig := `{
		"database": {
			"path": "test.db"
		},
		"server": {
			"port": 8080
		},
		"rss": {
			"update_period": 5,
			"feeds": [
				"http://example.com/rss"
			]
		}
	}`

	// Создаем временный файл
	tmpFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Записываем конфиг во временный файл
	if err := os.WriteFile(tmpFile.Name(), []byte(tmpConfig), 0644); err != nil {
		t.Fatal(err)
	}

	// Тестируем загрузку конфига
	cfg, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Проверяем значения
	if cfg.Database.Path != "test.db" {
		t.Errorf("Expected database path 'test.db', got '%s'", cfg.Database.Path)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", cfg.Server.Port)
	}
	if cfg.RSS.UpdatePeriod != 5 {
		t.Errorf("Expected update period 5, got %d", cfg.RSS.UpdatePeriod)
	}
	if len(cfg.RSS.Feeds) != 1 || cfg.RSS.Feeds[0] != "http://example.com/rss" {
		t.Errorf("Unexpected feeds configuration")
	}
}

func TestLoadConfigError(t *testing.T) {
	// Тестируем загрузку несуществующего файла
	_, err := LoadConfig("nonexistent.json")
	if err == nil {
		t.Error("Expected error when loading nonexistent file")
	}
}
