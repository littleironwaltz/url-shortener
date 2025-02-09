package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestURLShortener(t *testing.T) {
	store := NewURLStore()
	
	// Test /shorten endpoint
	t.Run("Shorten URL", func(t *testing.T) {
		reqBody := ShortenRequest{URL: "https://example.com"}
		body, err := json.Marshal(reqBody)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))
		w := httptest.NewRecorder()

		http.DefaultServeMux.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var resp ShortenResponse
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if !strings.HasPrefix(resp.ShortURL, "http://localhost:8080/") {
			t.Errorf("Expected short URL to start with http://localhost:8080/, got %s", resp.ShortURL)
		}

		code := strings.TrimPrefix(resp.ShortURL, "http://localhost:8080/")
		if len(code) != 6 {
			t.Errorf("Expected code length to be 6, got %d", len(code))
		}
	})

	// Test invalid request to /shorten
	t.Run("Invalid Shorten Request", func(t *testing.T) {
		reqBody := `{"url": ""}`
		req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		http.DefaultServeMux.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	// Test redirect endpoint
	t.Run("Redirect", func(t *testing.T) {
		// First, create a short URL
		originalURL := "https://example.com"
		code := "testcode"
		store.Set(code, originalURL)

		req := httptest.NewRequest(http.MethodGet, "/"+code, nil)
		w := httptest.NewRecorder()

		http.DefaultServeMux.ServeHTTP(w, req)

		if w.Code != http.StatusFound {
			t.Errorf("Expected status code %d, got %d", http.StatusFound, w.Code)
		}

		location := w.Header().Get("Location")
		if location != originalURL {
			t.Errorf("Expected redirect to %s, got %s", originalURL, location)
		}
	})

	// Test non-existent code
	t.Run("Non-existent Code", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
		w := httptest.NewRecorder()

		http.DefaultServeMux.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}
