package socket

import (
	"net/http"

	"socket-chatroom/core"

	"github.com/gorilla/websocket"
)

func ServerWS(hub core.Hubs, upgrader websocket.Upgrader, w http.ResponseWriter, r *http.Request) {
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
	room.Clients[client] = true

	go core.WritePump(client)
	hub.ReadPump(roomName, client)
}
