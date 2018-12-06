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
		connections: map[int64]*websocket.Conn{},
		mx:          sync.RWMutex{},
	}
	return n
}

//AddClient adds a new client to the Notifier
func (n *Notifier) AddClient(client *websocket.Conn, id int64) {
	log.Println("adding new WebSockets client")
	n.mx.Lock()
	n.connections[id] = client
	log.Printf("this is n.con %v", n.connections[id])
	fmt.Println("The map is of length: %v", len(n.connections))
	n.mx.Unlock()
	//TODO: add the client to the slice you are using
	//to track all current WebSocket connections.
	//Since this can be called from multiple
	//goroutines, make sure you protect that slice
	//while you add a new connection to it!
	// go n.readLoop(client)
	//also process incoming control messages from
	//the client, as described in this section of the docs:
	//https://godoc.org/github.com/gorilla/websocket#hdr-Control_Messages
}

// func (n *Notifier) readLoop(c *websocket.Conn) {
// 	for {
// 		if messageType, p, err := c.NextReader(); err != nil {
// 			//print
// 			log.Println(messageType, " ", p, " ", err)
// 			c.Close()
// 			break
// 		}
// 	}
// }

// Task is the struct that stands for Task object from microservice
type Task struct {
	Description    string `json:"description,omitempty"`
	Point          int    `json:"point,omitempty"`
	IsProgress     bool   `json:"is_progress,omitempty"`
	FamilyRoomName string `json:"familyRoomName,omitempty"`
	UserID         int    `json:"user_id,omitempty"`
}

// Message stands for the message queue from the Task microservice
type Message struct {
	Name  string      `json:"name,omitempty"`
	Task  *Task       `json:"task,omitempty"`
	Tasks []*Task     `json:"tasks,omitempty"`
	Point int         `json:"point,omitempty"`
	User  *users.User `json:"user,omitempty"`
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
	// for msg := range msgs {
	// n.mx.Lock()
	// defer n.mx.Unlock()
	for {
		d := <-msgs
		log.Printf("Received a task: %v", string(d.Body[:]))
		m := &Message{}
		if err := json.Unmarshal(d.Body, m); err != nil {
			fmt.Errorf("Error while unmarshal of d.Body: %v", err)
			return
		}

		log.Printf("Debug: start messages: %v", m)
		// // Get the roomname using getbyroomname() from mysql
		users, err := ctx.User.GetByRoomName(m.Task.FamilyRoomName)
		if err != nil {
			fmt.Errorf("Error while running GetByRoomName: %v", err)
			return
		}
		log.Printf("Debug: start users: %v", users)
		log.Printf(m.Task.FamilyRoomName)
		// Get admin using getadmin() for adding admin to users
		admin, err := ctx.User.GetAdmin(m.Task.FamilyRoomName, "Admin")
		if err != nil {
			fmt.Errorf("Error while running GetAdmin: %v", err)
			return
		}
		log.Printf("Debug: start admin2: %v", admin)
		users = append(users, admin)
		// if the done
		if m.Name == "task-done" {
			// should update points to mysql
			if _, err := ctx.User.UpdateScore(m.User.ID, m.Point); err != nil {
				fmt.Errorf("Error while running UpdateScore: %v", err)
				return
			}
		}
		// for loop through family members and write the message to the connections
		for _, user := range users {
			conn := n.connections[user.ID]
			// if Writemessage has an error, break the loop.
			if err := conn.WriteMessage(websocket.TextMessage, d.Body); err != nil {
				n.RemoveConnection(conn, user.ID)
				break
			}
		}

	}
}
