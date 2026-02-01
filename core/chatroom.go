package core

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Rooms map[string]*Room
	mu    sync.RWMutex
}

type Room struct {
	Clients map[*Client]bool
	Mu      sync.RWMutex
}

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

type Response struct {
	Room    string `json:"room"`
	Message string `json:"message"`
}

type Hubs interface {
	GetRoom(name string) *Room
	Broadcast(msg []byte)
	ReadPump(client *Client)
	WritePump(client *Client)
}

func NewHubs() Hubs {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

func (h *Hub) GetRoom(name string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, ok := h.Rooms[name]
	if !ok {
		room = &Room{
			Clients: make(map[*Client]bool),
		}
		h.Rooms[name] = room
	}

	return room
}

func (h *Hub) Broadcast(msg []byte) {
	var data Response
	if err := json.Unmarshal(msg, &data); err != nil {
		log.Println("invalid json:", err)
		return
	}

	h.mu.RLock()
	room, ok := h.Rooms[data.Room]
	h.mu.RUnlock()
	if !ok {
		return
	}

	room.Mu.RLock()
	defer room.Mu.RUnlock()

	for client := range room.Clients {
		select {
		case client.Send <- msg:
		default:
			close(client.Send)
			delete(room.Clients, client)
		}
	}
}

func (h *Hub) ReadPump(client *Client) {
	defer client.Conn.Close()

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}
		h.Broadcast(msg)
	}
}

func (h *Hub) WritePump(client *Client) {
	defer client.Conn.Close()

	for msg := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}
