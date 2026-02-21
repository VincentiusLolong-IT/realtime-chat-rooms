package routers

import (
	"net/http"
	"socket-chatroom/interface/controllers"
)

type Router interface {
	NewRoutInnit() *http.ServeMux
}

type Routers struct {
	ct controllers.Controller
}

func NewRouters(controllers controllers.Controller) Router {
	return &Routers{
		ct: controllers,
	}
}

func (r *Routers) NewRoutInnit() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", r.ct.Roomchat)
	return mux
}
