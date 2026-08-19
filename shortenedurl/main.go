package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"
)

type URLStore struct {
	mu		sync.RWMutex
	data	map[string]string
}

func NewURLStore() *URLStore {
	return &URLStore{data: make(map[string]string)}
}

func (s *URLStore) Save(longURL string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	code := generateCode(6)
	s.data[code] = longURL
	return code
}

func (s *URLStore) Get(code string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, ok := s.data[code]
	return url, ok
}

func generateCode(n int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, n)
	for i := range code {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		code[i] = chars[idx.Int64()]
	}
	return string(code)
}

var store = NewURLStore()

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if req.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	code := store.Save(req.URL)
	resp := shortenResponse{
		ShortCode: code,
		ShortURL: fmt.Sprintf("http://localhost:8099/%s", code),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:]
	if code == "" {
		http.NotFound(w, r)
		return
	}

	longURL, ok := store.Get(code)
	if !ok {
		http.NotFound(w, r)
		return
	}

	fmt.Println("redirect: %s -> %s", code, longURL)
	http.Redirect(w, r, longURL, http.StatusFound)
}

func main() {
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)

	fmt.Println("listening on :8099")
	http.ListenAndServe(":8099", nil)
}


