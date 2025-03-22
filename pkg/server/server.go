package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"anews/pkg/repository"
)

type Server struct {
	port int
	repo *repository.NewsRepository
}

func NewServer(port int, repo *repository.NewsRepository) *Server {
	return &Server{
		port: port,
		repo: repo,
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/api/news/", s.handleGetNews)
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", s.handleHome)

	addr := fmt.Sprintf(":%d", s.port)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handleGetNews(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")

	var limit int
	var err error

	// Check if a limit is provided in the URL
	if len(parts) == 4 && parts[3] != "" {
		limit, err = strconv.Atoi(parts[3])
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	} else {
		// If no limit is provided, fetch all news
		limit = -1 // Use -1 to indicate no limit
	}

	news, err := s.repo.GetLatestNews(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}
