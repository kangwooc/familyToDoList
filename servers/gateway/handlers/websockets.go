package handlers

import (
	"final-project-zco/servers/gateway/sessions"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// return r.Header.Get("Origin") != "https://kangwoo.tech"
		return true
	},
}

//WebSocketsHandler is a handler for WebSocket upgrade requests
type WebSocketsHandler struct {
	notifier *Notifier
	upgrader *websocket.Upgrader
}

//NewWebSocketsHandler constructs a new WebSocketsHandler
func NewWebSocketsHandler(notifier *Notifier) *WebSocketsHandler {
	//create, initialize, and return a new WebSocketsHandler
	return &WebSocketsHandler{
		notifier: notifier,
		upgrader: &upgrader,
	}
}

// WebSocketsHandler implements the http.Handler interface for the WebSocketsHandler
func (ctx *HandlerContext) WebSocketsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Debug: received websocket upgrade request")
	sessionState := &SessionState{}
	// log.Printf("Debug: r: %v", r)
	_, err := sessions.GetState(r, ctx.SigningKey, ctx.Session, sessionState)
	if err != nil {
		log.Printf("Debug: Error: %v", err)
		http.Error(w, "User must be authenticated", http.StatusUnauthorized)
		return
	}
	if !upgrader.CheckOrigin(r) {
		http.Error(w, "Origin not allowed", http.StatusForbidden)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "can't upgrade connection", http.StatusInternalServerError)
		return
	}
	log.Println("adding client to notifier")
	ctx.Socket.notifier.AddClient(conn, sessionState.User.ID)
}
