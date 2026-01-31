package controllers

import (
	"net/http"
	"socket-chatroom/core"

	"github.com/gorilla/websocket"
)

type Controller interface {
	Roomchat(w http.ResponseWriter, r *http.Request)
}

type Controllers struct {
	Hub      core.Hubs
	Upgrader websocket.Upgrader
}

func NewControllers(hub core.Hubs, upgrader websocket.Upgrader) Controller {
	return &Controllers{
		Hub:      hub,
		Upgrader: upgrader,
	}
}
