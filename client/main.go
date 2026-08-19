package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)


type Post struct {
	ID		int 	`json:"id"`
	Title	string	`json:"title"`
}

func main() {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var p Post
	if err := json.Unmarshal(body, &p); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Title: %s\n", p.Title)
}

