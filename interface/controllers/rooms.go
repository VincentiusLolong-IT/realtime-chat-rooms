package controllers

import (
	"net/http"

	sock "socket-chatroom/interface/controllers/socket"
)

func (c *Controllers) Roomchat(w http.ResponseWriter, r *http.Request) {
	sock.ServerWS(c.Hub, c.Upgrader, w, r)
}
