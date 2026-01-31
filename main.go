package main

import (
	"net/http"
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
	send chan []byte
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func newHub() *Hub {
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

func (h *Hub) Broadcast(roomName string, msg []byte) {
	h.mu.Lock()
	room, ok := h.Rooms[roomName]
	h.mu.Unlock()

	if !ok {
		return
	}

	for client := range room.Clients {
		select {
		case client.send <- msg:
		default:
			close(client.send)
			delete(room.Clients, client)
		}
	}
}

func writePump(client *Client) {
	defer client.Conn.Close()

	for msg := range client.send {
		err := client.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}

func readPump(hub *Hub, roomName string, client *Client) {
	defer client.Conn.Close()

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}
		hub.Broadcast(roomName, msg)
	}
}

func serverWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	roomName := r.URL.Query().Get("room")
	if roomName == "" {
		roomName = "general"
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		Conn: conn,
		send: make(chan []byte, 256),
	}

	room := hub.GetRoom(roomName)
	room.Clients[client] = true

	go writePump(client)
	readPump(hub, roomName, client)
}

func main() {
	hub := newHub()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serverWS(hub, w, r)
	})

	http.ListenAndServe(":8080", nil)
}
