package models

import "time"

type News struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PubDate     time.Time `json:"pub_date"`
	SourceURL   string    `json:"source_url"`
	CreatedAt   time.Time `json:"created_at"`
}
