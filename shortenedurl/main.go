package main

import (
	"crypto/rand"
	"encoding/json"
	"database/sql"
	"fmt"
	"math/big"
	"net/http"
)

var db *sql.DB

func generateCode(n int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, n)
	for i := range code {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		code[i] = chars[idx.Int64()]
	}
	return string(code)
}

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

	code := generateCode(6)
	_, err := db.Exec(`INSERT INTO urls (code, long_url) VALUES ($1, $2)`, code, req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := shortenResponse {
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

	var longURL string
	err := db.QueryRow(`SELECT long_url FROM urls WHERE code = $1`, code).Scan(&longURL)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

func main() {
	db = NewDB()
	defer db.Close()

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)

	fmt.Println("listening on :8099")
	http.ListenAndServe(":8099", nil)
}


