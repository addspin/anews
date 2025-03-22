package repository

import (
	"os"
	"testing"
	"time"

	"anews/pkg/models"
)

func TestNewsRepository(t *testing.T) {
	// Создаем временную БД для тестов
	tmpDB := "test.db"
	defer os.Remove(tmpDB)

	repo, err := NewNewsRepository(tmpDB)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.db.Close()

	// Тестируем сохранение новости
	testNews := &models.News{
		Title:       "Test Title",
		Description: "Test Description",
		PubDate:     time.Now().UTC(),
		SourceURL:   "http://example.com/test",
	}

	if err := repo.SaveNews(testNews); err != nil {
		t.Fatalf("Failed to save news: %v", err)
	}

	// Проверяем, что ID был установлен
	if testNews.ID == 0 {
		t.Error("Expected ID to be set after save")
	}

	// Тестируем получение новостей
	news, err := repo.GetLatestNews(10)
	if err != nil {
		t.Fatalf("Failed to get news: %v", err)
	}

	if len(news) != 1 {
		t.Fatalf("Expected 1 news item, got %d", len(news))
	}

	// Проверяем поля первой новости
	if news[0].Title != testNews.Title {
		t.Errorf("Expected title %s, got %s", testNews.Title, news[0].Title)
	}

	// Тестируем уникальность URL
	dupNews := &models.News{
		Title:       "Another Title",
		Description: "Another Description",
		PubDate:     time.Now().UTC(),
		SourceURL:   "http://example.com/test", // тот же URL
	}

	if err := repo.SaveNews(dupNews); err != nil {
		t.Fatalf("Failed to handle duplicate URL: %v", err)
	}

	// Проверяем, что дубликат не был добавлен
	news, _ = repo.GetLatestNews(10)
	if len(news) != 1 {
		t.Errorf("Expected still 1 news item after duplicate, got %d", len(news))
	}
}
