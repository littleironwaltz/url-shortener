package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type URLStore struct {
	mu    sync.RWMutex
	store map[string]string
}

func NewURLStore() *URLStore {
	return &URLStore{
		store: make(map[string]string),
	}
}

func (s *URLStore) Set(code, url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[code] = url
}

func (s *URLStore) Get(code string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, exists := s.store[code]
	return url, exists
}

func generateCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 6
	
	code := make([]byte, codeLength)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

func setupHandlers(store *URLStore) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ShortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.URL == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		code := generateCode()
		store.Set(code, req.URL)

		resp := ShortenResponse{
			ShortURL: fmt.Sprintf("http://localhost:8080/%s", code),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		code := r.URL.Path[1:]
		if code == "" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		url, exists := store.Get(code)
		if !exists {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, url, http.StatusFound)
	})

	return mux
}

func main() {
	rand.Seed(time.Now().UnixNano())
	store := NewURLStore()
	
	handler := setupHandlers(store)
	
	log.Printf("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
