package repository

import (
	"database/sql"
	"time"

	"anews/pkg/models"

	_ "github.com/mattn/go-sqlite3"
)

type NewsRepository struct {
	db *sql.DB
}

func NewNewsRepository(dbPath string) (*NewsRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Создаем таблицу, если она не существует
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS news (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			pub_date DATETIME NOT NULL,
			source_url TEXT NOT NULL UNIQUE,
			created_at DATETIME NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return &NewsRepository{db: db}, nil
}

func (r *NewsRepository) GetLatestNews(limit int) ([]models.News, error) {
	query := `
		SELECT id, title, description, datetime(pub_date), source_url, datetime(created_at)
		FROM news 
		ORDER BY pub_date DESC 
	`

	var rows *sql.Rows
	var err error

	if limit < 0 {
		// No limit, fetch all news
		rows, err = r.db.Query(query)
	} else {
		// Apply limit
		query += "LIMIT ?"
		rows, err = r.db.Query(query, limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var news []models.News
	for rows.Next() {
		var n models.News
		var pubDate, createdAt string
		err := rows.Scan(&n.ID, &n.Title, &n.Description, &pubDate, &n.SourceURL, &createdAt)
		if err != nil {
			return nil, err
		}
		n.PubDate, _ = time.Parse("2006-01-02 15:04:05", pubDate)
		n.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		news = append(news, n)
	}
	return news, nil
}

func (r *NewsRepository) SaveNews(news *models.News) error {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM news WHERE source_url = ?)", news.SourceURL).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	query := `INSERT INTO news (title, description, pub_date, source_url, created_at)
		VALUES (?, ?, datetime(?), ?, datetime(?))`

	result, err := r.db.Exec(
		query,
		news.Title,
		news.Description,
		news.PubDate.UTC().Format("2006-01-02 15:04:05"),
		news.SourceURL,
		time.Now().UTC().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	news.ID = id
	return nil
}
