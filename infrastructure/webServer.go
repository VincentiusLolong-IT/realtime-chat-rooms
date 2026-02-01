package infrastructure

import (
	"log"
	"net/http"
	"socket-chatroom/core"
	"socket-chatroom/infrastructure/routers"
	"socket-chatroom/interface/controllers"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebServices() {
	newHub := core.NewHubs()
	newController := controllers.NewControllers(newHub, upgrader)
	newRouters := routers.NewRouters(newController)

	mux := newRouters.NewRoutInnit()

	log.Fatal(http.ListenAndServe(":8080", mux))
}
