package rss

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"anews/pkg/models"
)

type Feed struct {
	Channel struct {
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Items       []Item `xml:"item"`
	} `xml:"channel"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Link        string `xml:"link"`
	Date        string `xml:"date"`
	Published   string `xml:"published"`
	Updated     string `xml:"updated"`
}

type Reader struct {
	feeds []string
}

func NewReader(feeds []string) *Reader {
	return &Reader{feeds: feeds}
}

func (r *Reader) ReadFeed(feed string) ([]models.News, error) {
	log.Printf("Fetching feed: %s", feed)
	resp, err := http.Get(feed)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rssFeed Feed
	if err := xml.NewDecoder(resp.Body).Decode(&rssFeed); err != nil {
		return nil, err
	}

	var news []models.News
	for _, item := range rssFeed.Channel.Items {
		pubDateStr := item.PubDate
		if pubDateStr == "" {
			pubDateStr = item.Date
		}
		if pubDateStr == "" {
			pubDateStr = item.Published
		}
		if pubDateStr == "" {
			pubDateStr = item.Updated
		}

		pubDate, err := parseRSSDate(pubDateStr)
		if err != nil {
			log.Printf("Warning: cannot parse date '%s': %v", pubDateStr, err)
			continue
		}

		news = append(news, models.News{
			Title:       item.Title,
			Description: item.Description,
			PubDate:     pubDate,
			SourceURL:   item.Link,
		})
	}

	return news, nil
}

func parseRSSDate(date string) (time.Time, error) {
	date = strings.TrimSpace(date)
	if strings.Contains(date, "GMT") {
		date = strings.Replace(date, "GMT", "+0000", 1)
	}

	formats := []string{
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"Mon, 02 Jan 2006 15:04:05 +0000",
		"Mon, 02 Jan 2006 15:04:05 MST",
		"Mon, 02 Jan 2006 15:04:05 Z",
		"Mon, 02 Jan 2006 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		t, err := time.Parse(format, date)
		if err == nil {
			if t.Location() == time.Local {
				t = t.UTC()
			}
			return t, nil
		}
	}

	if ts, err := strconv.ParseInt(date, 10, 64); err == nil {
		return time.Unix(ts, 0), nil
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", date)
}

// ... остальные методы из internal/rss/reader.go ...
