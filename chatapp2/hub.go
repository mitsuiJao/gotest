package main

import (
	"context"
	"sync"
)

type Hub struct {
	rooms map[string]*Room
	mu    sync.Mutex
	ctx   context.Context
}

func newHub(ctx context.Context) *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
		ctx:   ctx,
	}
}

func (h *Hub) getOrCreateRoom(roomID string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.rooms[roomID]; ok {
		return room
	}

	room := newRoom(h.ctx, roomID)
	h.rooms[roomID] = room
	go room.run()
	go func() {
		<-room.ctx.Done()
		h.mu.Lock()
		delete(h.rooms, roomID)
		h.mu.Unlock()
	}()

	return room
}
