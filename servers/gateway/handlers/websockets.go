package handlers

import (
	"final-project-zco/servers/gateway/sessions"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// //NotificationsHandler handles requests for the /notifications resource
// type NotificationsHandler struct {
// 	notifier *Notifier
// }

// //NewNotificationsHandler constructs a new NotificationsHandler
// func NewNotificationsHandler(notifier *Notifier) *NotificationsHandler {
// 	return &NotificationsHandler{notifier}
// }

// //ServeHTTP handles HTTP requests for the NotificationsHandler
// func (nh *NotificationsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	//NOTE: this is just a simple handler for testing that
// 	//triggers a new notification anytime this handler
// 	//receives an HTTP request using any method.
// 	//In your real server, you will listen for new messages
// 	//from your MQ server, and pass them to the Notifier as
// 	//you receive them.
// 	w.Header().Add("Access-Control-Allow-Origin", "*")
// 	msg := fmt.Sprintf("Notification pushed from the server at %s", time.Now().Format("15:04:05"))
// 	nh.notifier.Notify([]byte(msg))
// }

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
	//TODO: add a field for the websocket.Upgrader
	//see https://godoc.org/github.com/gorilla/websocket
	//and https://godoc.org/github.com/gorilla/websocket#Upgrader
}

//NewWebSocketsHandler constructs a new WebSocketsHandler
func NewWebSocketsHandler(notifier *Notifier) *WebSocketsHandler {
	//create, initialize, and return a new WebSocketsHandler
	return &WebSocketsHandler{
		notifier: notifier,
		upgrader: &upgrader,
	}
}

//ServeHTTP implements the http.Handler interface for the WebSocketsHandler
func (ctx *HandlerContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("received websocket upgrade request")
	sessionState := &SessionState{}
	_, err := sessions.GetState(r, ctx.SigningKey, ctx.Session, sessionState)
	if err != nil {
		http.Error(w, "User must be authenticated", http.StatusUnauthorized)
		return
	}

	if !upgrader.CheckOrigin(r) {
		http.Error(w, "Origin not allowed", http.StatusForbidden)
		return
	}
	//TODO: upgrade the connection to a WebScoket
	//see https://godoc.org/github.com/gorilla/websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("adding client to notifier")
	//TODO: add the new WebSocket connection to the Notifier
	ctx.Socket.notifier.AddClient(conn, sessionState.User.ID)
}
