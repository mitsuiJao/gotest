package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade failed: ", err)
		return
	}
	defer conn.Close()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error: ", err)
			break
		}
		log.Printf("received: %s\n", msg)

		if err := conn.WriteMessage(msgType, msg); err != nil {
			log.Println("write error: ", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", echoHandler)

	fmt.Println("listening on :8095")
	log.Fatal(http.ListenAndServe(":8095", nil))
}
