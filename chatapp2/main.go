package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade failed:", err)
		return
	}

	conn.WriteMessage(websocket.TextMessage, []byte("room id: "))
	_, msg1, err := conn.ReadMessage()
	if err != nil {
		conn.Close()
		return
	}
	roomID := string(msg1)

	room := hub.getOrCreateRoom(roomID)

	conn.WriteMessage(websocket.TextMessage, []byte("yourname: "))
	_, msg2, err := conn.ReadMessage()
	if err != nil {
		conn.Close()
		return
	}

	client := &Client{
		name: string(msg2),
		room: room,
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.room.register <- client

	go client.writePump()
	go client.readPump()
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	hub := newHub(ctx)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	fmt.Println("listening on :8096")
	srv := &http.Server{
		Addr: ":8096",
	}
	go srv.ListenAndServe()
	<-ctx.Done()

	srv.Shutdown(context.Background())
}
