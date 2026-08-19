package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	ID		int		`json:"id"`
	Name	string	`json:"name"`
}

var users =[]User {
	{ID: 1, Name: "Taro"},
	{ID: 2, Name: "Hanako"},
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(users)
	case http.MethodPost:
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u.ID = len(users) + 1
		users = append(users, u)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
	default:
		http.Error(w, "method not allowd", http.StatusMethodNotAllowed)
	}
}


func main() {
	http.HandleFunc("/users", userHandler)
	log.Println("listening on :55555")
	log.Fatal(http.ListenAndServe(":55555", nil))
}

