package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"final-project-zco/servers/gateway/models/users"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

//Notifier is an object that handles WebSocket notifications.
//It starts a new
type Notifier struct {
	// eventQ chan []byte
	// map for updating the task list
	connections map[int64]*websocket.Conn
	mx          sync.RWMutex
}

//NewNotifier constructs a new Notifier
func NewNotifier() *Notifier {
	n := &Notifier{
		connections: make(map[int64]*websocket.Conn),
		mx:          sync.RWMutex{},
	}
	return n
}

//AddClient adds a new client to the Notifier
func (n *Notifier) AddClient(client *websocket.Conn, id int64) {
	log.Println("adding new WebSockets client")
	n.mx.Lock()
	defer n.mx.Unlock()
	n.connections[id] = client
	//TODO: add the client to the slice you are using
	//to track all current WebSocket connections.
	//Since this can be called from multiple
	//goroutines, make sure you protect that slice
	//while you add a new connection to it!
	go n.readLoop(client)
	//also process incoming control messages from
	//the client, as described in this section of the docs:
	//https://godoc.org/github.com/gorilla/websocket#hdr-Control_Messages
}
func (n *Notifier) readLoop(c *websocket.Conn) {
	for {
		if _, _, err := c.NextReader(); err != nil {
			c.Close()
			break
		}
	}
}

// Task is the struct that stands for Task object from microservice
type Task struct {
	description    string
	point          int
	isProgress     bool
	familyID       int
	familyRoomName string
	userID         int
}

// Message stands for the message queue from the Task microservice
type Message struct {
	name  string
	task  *Task
	tasks []*Task
	point int
	user  *users.User
}

// RemoveConnection is the method that remove connection
// It is a thread-safe method for inserting a connection
func (n *Notifier) RemoveConnection(conn *websocket.Conn, userID int64) {
	n.mx.Lock()
	// delete socket connection
	delete(n.connections, userID)
	n.mx.Unlock()
}

// Start starts the notification loop
func (n *Notifier) Start(msgs <-chan amqp.Delivery, name string, ctx *HandlerContext) {
	log.Println("starting notifier loop")
	// for msg := range msgs {
	n.mx.Lock()
	defer n.mx.Unlock()
	for {
		for d := range msgs {
			// if name == "taskQueue" {
			m := &Message{}
			if err := json.Unmarshal(d.Body, m); err != nil {
				fmt.Errorf("Error while unmarshal of d.Body: %v", err)
				return
			}
			// Get the roomname using getbyroomname() from mysql
			users, err := ctx.User.GetByRoomName(m.task.familyRoomName)
			if err != nil {
				fmt.Errorf("Error while running GetByRoomName: %v", err)
				return
			}
			// Get admin using getadmin() for adding admin to users
			admin, err := ctx.User.GetAdmin(m.task.familyRoomName, "Admin")
			if err != nil {
				fmt.Errorf("Error while running GetAdmin: %v", err)
				return
			}
			users = append(users, admin)
			// if the done
			if m.name == "task-done" {
				// should update points to mysql
				if _, err := ctx.User.UpdateScore(m.user.ID, m.point); err != nil {
					fmt.Errorf("Error while running UpdateScore: %v", err)
					return
				}
			}
			// for loop through family members and write the message to the connections
			for _, user := range users {
				conn := n.connections[user.ID]
				// if Writemessage has an error, break the loop.
				if err := conn.WriteMessage(1, d.Body); err != nil {
					n.RemoveConnection(conn, user.ID)
					break
				}
			}
			// } else {

			// }
		}
	}

	//TODO: start a never-ending loop that reads
	//new events out of the `n.eventQ` and broadcasts
	//them to all WebSocket connections.
	//To write the byte-slice to the WebSocket, use
	//the .WriteMessage() method.
	//https://godoc.org/github.com/gorilla/websocket#Conn.WriteMessage
	//Or, for better performance, prepare the message once
	//and use the .WritePreparedMessage() method.
	//https://godoc.org/github.com/gorilla/websocket#PreparedMessage

	//Remember that you need to lock the slice of connections
	//while you iterate it, as other goroutines might
	//be trying to add new clients to it while you iterate!

	//If you get an error while trying to write the
	//message to one of the WebSocket connections,
	//that means the client has disconnected, so
	//remove that connection from your list.
}
