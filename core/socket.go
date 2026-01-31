package core

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Rooms map[string]*Room
	mu    sync.RWMutex
}

type Room struct {
	Clients map[*Client]bool
}

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

type Hubs interface {
	GetRoom(name string) *Room
	Broadcast(roomName string, msg []byte)
	ReadPump(roomName string, client *Client)
}

func NewHubs() Hubs {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

func (h *Hub) GetRoom(name string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.Rooms == nil {
		h.Rooms = make(map[string]*Room)
	}

	room, ok := h.Rooms[name]
	if !ok {
		room = &Room{
			Clients: make(map[*Client]bool),
		}
		h.Rooms[name] = room
	}

	return room
}

func (h *Hub) Broadcast(roomName string, msg []byte) {
	h.mu.Lock()
	room, ok := h.Rooms[roomName]
	h.mu.Unlock()

	if !ok {
		return
	}

	for client := range room.Clients {
		select {
		case client.Send <- msg:
		default:
			close(client.Send)
			delete(room.Clients, client)
		}
	}
}

func WritePump(client *Client) {
	defer client.Conn.Close()

	for msg := range client.Send {
		err := client.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}

func (h *Hub) ReadPump(roomName string, client *Client) {
	defer client.Conn.Close()

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}
		h.Broadcast(roomName, msg)
	}
}
