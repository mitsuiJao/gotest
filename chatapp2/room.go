package main

import (
	"context"
	"log"
)

type Room struct {
	id         string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	ctx        context.Context
	cancel     context.CancelFunc
}

func newRoom(parent context.Context, id string) *Room {
	ctx, cancel := context.WithCancel(parent)
	return &Room{
		id:         id,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		ctx:        ctx,
		cancel:     cancel,
	}
}

func disconnect(r *Room, c *Client) {
	if _, ok := r.clients[c]; ok {
		delete(r.clients, c)
		close(c.send)
		if len(r.clients) == 0 {
			r.cancel()
		}
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
