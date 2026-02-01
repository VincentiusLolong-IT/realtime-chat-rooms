package socket

import (
	"net/http"

	"socket-chatroom/core"

	"github.com/gorilla/websocket"
)

func ServerWS(
	hub core.Hubs,
	upgrader websocket.Upgrader,
	w http.ResponseWriter,
	r *http.Request,
) {
	roomName := r.URL.Query().Get("room")
	if roomName == "" {
		roomName = "general"
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &core.Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	room := hub.GetRoom(roomName)

	room.Mu.Lock()
	room.Clients[client] = true
	room.Mu.Unlock()

	defer func() {
		room.Mu.Lock()
		delete(room.Clients, client)
		room.Mu.Unlock()
		close(client.Send)
	}()

	go hub.WritePump(client)
	hub.ReadPump(client)
}
