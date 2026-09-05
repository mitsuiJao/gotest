package main

import (
	"context"
	"log"
)

type Room struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	ctx        context.Context
}

func newRoom(ctx context.Context) *Room {
	return &Room{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		ctx:        ctx,
	}
}

func disconnect(r *Room, c *Client) {
	if _, ok := r.clients[c]; ok {
		delete(r.clients, c)
		close(c.send)
	}
}

func (r *Room) run() {
	for {
		select {
		case <-r.ctx.Done():
			for c := range r.clients {
				disconnect(r, c)
			}
			log.Println("stop")
			return

		case client := <-r.register:
			r.clients[client] = true
			log.Printf("client registerd (total: %d)\n", len(r.clients))

		case client := <-r.unregister:
			disconnect(r, client)
			log.Printf("client unregistered (total: %d)", len(r.clients))

		case msg := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		}
	}
}
