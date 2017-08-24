package test

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/onsi/gomega/ghttp"
)

type WebSocketHandler func(*websocket.Conn)

func RouteToWSHandler(server *ghttp.Server, method, path string, handler WebSocketHandler) {
	server.RouteToHandler(method, path, wshandler(handler))
}

func WebSocketURL(server *ghttp.Server) *url.URL {
	u, _ := url.Parse(server.URL())
	u.Scheme = "ws"
	return u
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(h WebSocketHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Cannot upgrade: %v", err), http.StatusInternalServerError)
		}
		h(conn)
	}
}
