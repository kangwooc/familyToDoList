package handlers

import (
	"final-project-zco/servers/gateway/sessions"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

//TODO: add a handler that upgrades clients to a WebSocket connection
//and adds that to a list of WebSockets to notify when events are
//read from the RabbitMQ server. Remember to synchronize changes
//to this list, as handlers are called concurrently from multiple
//goroutines.

//TODO: start a goroutine that connects to the RabbitMQ server,
//reads events off the queue, and broadcasts them to all of
//the existing WebSocket connections that should hear about
//that event. If you get an error writing to the WebSocket,
//just close it and remove it from the list
//(client went away without closing from
//their end). Also make sure you start a read pump that
//reads incoming control messages, as described in the
//Gorilla WebSocket API documentation:
//http://godoc.org/github.com/gorilla/websocket

// SocketStore stores all the connections
type SocketStore struct {
	Connections map[int64][]*websocket.Conn
	lock        sync.RWMutex
}

// Control messages for websocket
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

func NewSocketStore() *SocketStore {
	return &SocketStore{make(map[int64][]*websocket.Conn), sync.RWMutex{}}
}

// Thread-safe method for inserting a connection
func (s *SocketStore) InsertConnection(conn *websocket.Conn, userid int64) int64 {
	s.lock.Lock()
	// insert socket connection
	s.Connections[userid] = append(s.Connections[userid], conn)
	s.lock.Unlock()
	return userid
}

// Thread-safe method for inserting a connection
func (s *SocketStore) RemoveConnection(conn *websocket.Conn, userID int64) {
	s.lock.Lock()
	// insert socket connection
	for k, v := range s.Connections[userID] {
		if v == conn {
			s.Connections[userID] = append(s.Connections[userID][:k], s.Connections[userID][k+1:]...)
		}
	}
	s.lock.Unlock()
}

// Simple method for writing a message to all live connections.
// In your homework, you will be writing a message to a subset of connections
// (if the message is intended for a private channel), or to all of them (if the message
// is posted on a public channel
func (s *SocketStore) WriteToConnections(m channel) error {
	s.lock.Lock()
	if len(m.UserIDs) == 0 { // public
		for i := 0; i < len(s.Connections); i++ {
			for _, k := range s.Connections[int64(i)] {
				messageType, p, err := k.ReadMessage()
				if err != nil {
					log.Println(err)
					return err
				}
				if err := k.WriteMessage(messageType, p); err != nil {
					log.Println(err)
					return err
				}
			}
		}
	} else { //private
		for i := 0; i < len(s.Connections); i++ {
			for _, l := range m.UserIDs {
				for _, k := range s.Connections[l] {
					messageType, p, err := k.ReadMessage()
					if err != nil {
						log.Println(err)
						return err
					}
					if err := k.WriteMessage(messageType, p); err != nil {
						log.Println(err)
						return err
					}
				}
			}

		}
	}
	return nil
}

// This is a struct to read our message into
type channel struct {
	Type    string      `json:"type"`
	Channel interface{} `json:"channel"`
	UserIDs []int64     `json:"userIDs"`
	Message interface{} `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// This function's purpose is to reject websocket upgrade requests if the
		// origin of the websockete handshake request is coming from unknown domains.
		// This prevents some random domain from opening up a socket with your server.
		return r.Header.Get("Origin") != "https://kangwooc.tech"
	},
}

func (ctx *HandlerContext) WebSocketConnectionHandler(w http.ResponseWriter, r *http.Request) {
	ss := &SessionState{}
	if _, err := sessions.GetState(r, ctx.SigningKey, ctx.Session, ss); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// handle the websocket handshake
	if !upgrader.CheckOrigin(r) {
		http.Error(w, "Websocket Connection Refused", http.StatusUnauthorized)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to open websocket connection", http.StatusUnauthorized)
		return
	}
	sockets := NewSocketStore()
	// Insert our connection onto our datastructure for ongoing usage
	userID := sockets.InsertConnection(conn, ss.User.ID)

	// forever := make(chan bool)
	// check error log
	go (func(conn *websocket.Conn, userID int64) {
		defer conn.Close()
		defer sockets.RemoveConnection(conn, userID)
		for { // infinite loop
			m := channel{}
			err := conn.ReadJSON(&m)
			if err != nil {
				fmt.Println("Error reading json.", err)
				conn.Close()
				break
			}
			// newMsg <- sharedChannel
			err = sockets.WriteToConnections(m)
			if err != nil {
				fmt.Println("Error writing to connection.", err)
				conn.Close()
				break
			}
			fmt.Printf("Got message: %#v\n", m)

			if err = conn.WriteJSON(m); err != nil {
				fmt.Println(err)
			}

		}
	})(conn, userID)

	// cleanup
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
