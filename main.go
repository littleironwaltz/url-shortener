// Package main implements a URL shortener service that provides functionality
// to create short URLs and redirect to their original destinations.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

// Logger represents a structured logger with different log levels
type Logger struct {
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
}

// NewLogger creates a new structured logger with different log levels
func NewLogger() *Logger {
	return &Logger{
		info:  log.New(os.Stdout, "INFO: ", log.LstdFlags),
		warn:  log.New(os.Stdout, "WARN: ", log.LstdFlags),
		error: log.New(os.Stderr, "ERROR: ", log.LstdFlags),
	}
}

// URLStore provides thread-safe storage for URL mappings
type URLStore struct {
	mu    sync.RWMutex
	store map[string]string
	log   *Logger
}

// NewURLStore creates a new URLStore instance with initialized storage and logger
func NewURLStore() *URLStore {
	return &URLStore{
		store: make(map[string]string),
		log:   NewLogger(),
	}
}

// Set stores a URL with the given code in a thread-safe manner
func (s *URLStore) Set(ctx context.Context, code, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	select {
	case <-ctx.Done():
		s.log.error.Printf("Context cancelled while setting URL for code %s", code)
		return ctx.Err()
	default:
		s.store[code] = url
		s.log.info.Printf("Stored URL for code %s", code)
		return nil
	}
}

// Get retrieves a URL by its code in a thread-safe manner
func (s *URLStore) Get(ctx context.Context, code string) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	select {
	case <-ctx.Done():
		s.log.error.Printf("Context cancelled while getting URL for code %s", code)
		return "", false, ctx.Err()
	default:
		url, exists := s.store[code]
		if !exists {
			s.log.warn.Printf("URL not found for code %s", code)
		} else {
			s.log.info.Printf("Retrieved URL for code %s", code)
		}
		return url, exists, nil
	}
}

// generateCode creates a random 6-character code for URL shortening
func generateCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 6
	
	code := make([]byte, codeLength)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// ShortenRequest represents the request body for URL shortening
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse represents the response body for URL shortening
type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

// setupHandlers configures and returns an HTTP handler with URL shortening endpoints
func setupHandlers(store *URLStore) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			store.log.warn.Printf("Method %s not allowed for /shorten", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ShortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			store.log.error.Printf("Failed to decode request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.URL == "" {
			store.log.warn.Printf("Empty URL provided in request")
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		code := generateCode()
		if err := store.Set(r.Context(), code, req.URL); err != nil {
			store.log.error.Printf("Failed to store URL: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		resp := ShortenResponse{
			ShortURL: fmt.Sprintf("http://localhost:8080/%s", code),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			store.log.error.Printf("Failed to encode response: %v", err)
			return
		}
		store.log.info.Printf("Successfully shortened URL with code %s", code)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			store.log.warn.Printf("Method %s not allowed for redirect", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		code := r.URL.Path[1:]
		if code == "" {
			store.log.warn.Printf("Empty code provided in request")
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		url, exists, err := store.Get(r.Context(), code)
		if err != nil {
			store.log.error.Printf("Failed to retrieve URL: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if !exists {
			store.log.warn.Printf("URL not found for code: %s", code)
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		store.log.info.Printf("Redirecting code %s to URL %s", code, url)
		http.Redirect(w, r, url, http.StatusFound)
	})

	return mux
}

func main() {
	rand.Seed(time.Now().UnixNano())
	store := NewURLStore()
	
	handler := setupHandlers(store)
	
	store.log.info.Printf("Starting server on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		store.log.error.Fatalf("Server failed to start: %v", err)
	}
}
