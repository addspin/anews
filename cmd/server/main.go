package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"anews/pkg/config"
	"anews/pkg/repository"
	"anews/pkg/rss"
	"anews/pkg/server"
)

func main() {
	configPath := flag.String("config", "config.json", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	repo, err := repository.NewNewsRepository(cfg.Database.Path)
	if err != nil {
		log.Fatal(err)
	}

	reader := rss.NewReader(cfg.RSS.Feeds)

	// Запускаем горутины для чтения RSS-лент
	var wg sync.WaitGroup
	for _, feed := range cfg.RSS.Feeds {
		wg.Add(1)
		go func(feed string) {
			defer wg.Done()
			for {
				news, err := reader.ReadFeed(feed)
				if err != nil {
					log.Printf("Error reading feed %s: %v", feed, err)
					continue
				}

				for _, item := range news {
					if err := repo.SaveNews(&item); err != nil {
						log.Printf("Error saving news: %v", err)
					}
				}

				time.Sleep(time.Duration(cfg.RSS.UpdatePeriod) * time.Minute)
			}
		}(feed)
	}

	// Запускаем HTTP-сервер
	srv := server.NewServer(cfg.Server.Port, repo)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
