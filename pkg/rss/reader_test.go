package rss

import (
	"testing"
	"time"
)

func TestRSSReader(t *testing.T) {
	// Тест с реальным RSS-фидом
	reader := NewReader([]string{
		"http://static.feed.rbc.ru/rbc/logical/footer/news.rss",
	})

	news, err := reader.ReadFeed("http://static.feed.rbc.ru/rbc/logical/footer/news.rss")
	if err != nil {
		t.Fatalf("Failed to read RSS feed: %v", err)
	}

	if len(news) == 0 {
		t.Fatal("Expected some news items, got none")
	}

	// Проверяем первую новость
	firstNews := news[0]
	if firstNews.Title == "" {
		t.Error("Expected non-empty title")
	}
	if firstNews.Description == "" {
		t.Error("Expected non-empty description")
	}
	if firstNews.SourceURL == "" {
		t.Error("Expected non-empty source URL")
	}
	if firstNews.PubDate.IsZero() {
		t.Error("Expected non-zero publication date")
	}
	if firstNews.PubDate.After(time.Now()) {
		t.Error("Publication date should not be in the future")
	}
}

func TestParseRSSDate(t *testing.T) {
	tests := []struct {
		input     string
		wantErr   bool
		checkYear bool
	}{
		{"Mon, 02 Jan 2006 15:04:05 GMT", false, true},
		{"2006-01-02T15:04:05Z", false, true},
		{"invalid date", true, false},
		{"Mon, 02 Jan 2006 15:04:05 +0000", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			date, err := parseRSSDate(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if tt.checkYear && date.Year() < 2000 {
				t.Errorf("Expected year >= 2000, got %d", date.Year())
			}
		})
	}
}
