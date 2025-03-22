package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewsJSON(t *testing.T) {
	now := time.Now().UTC()
	news := News{
		ID:          1,
		Title:       "Test Title",
		Description: "Test Description",
		PubDate:     now,
		SourceURL:   "http://example.com",
		CreatedAt:   now,
	}

	// Тестируем маршалинг в JSON
	data, err := json.Marshal(news)
	if err != nil {
		t.Fatalf("Failed to marshal news: %v", err)
	}

	// Тестируем анмаршалинг из JSON
	var decoded News
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal news: %v", err)
	}

	// Проверяем поля
	if decoded.ID != news.ID {
		t.Errorf("Expected ID %d, got %d", news.ID, decoded.ID)
	}
	if decoded.Title != news.Title {
		t.Errorf("Expected Title %s, got %s", news.Title, decoded.Title)
	}
	if decoded.Description != news.Description {
		t.Errorf("Expected Description %s, got %s", news.Description, decoded.Description)
	}
	if decoded.SourceURL != news.SourceURL {
		t.Errorf("Expected SourceURL %s, got %s", news.SourceURL, decoded.SourceURL)
	}
	if !decoded.PubDate.Equal(news.PubDate) {
		t.Errorf("Expected PubDate %v, got %v", news.PubDate, decoded.PubDate)
	}
	if !decoded.CreatedAt.Equal(news.CreatedAt) {
		t.Errorf("Expected CreatedAt %v, got %v", news.CreatedAt, decoded.CreatedAt)
	}
}
