package test

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega/ghttp"
)

// WebSocketHandler is a generic WebSocket handler that a user would provide in
// a test
type WebSocketHandler func(*websocket.Conn)

// RouteToWSHandler is called on a ghttp.Server to set up a route to the
// provided WebSocketHandler.
func RouteToWSHandler(server *ghttp.Server, method, path string, handler WebSocketHandler) {
	server.RouteToHandler(method, path, wshandler(handler))
}

// WebSocketURL ensures that the URL passed uses the websocket scheme.
func WebSocketURL(server *ghttp.Server) *url.URL {
	u, _ := url.Parse(server.URL())
	u.Scheme = "ws"
	return u
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// wshandler is a trivial wrapper around a WebSocketHandler that turns it into
// an ordinary http.HandlerFunc.
func wshandler(h WebSocketHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			ginkgo.Fail(fmt.Sprintf("Cannot upgrade: %v", err))
		}
		h(conn)
	}
}
