package handlers

import (
	"log"
	"sync"

	"final-project-zco/servers/gateway/models/users"

	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
)

//Notifier is an object that handles WebSocket notifications.
//It starts a new
type Notifier struct {
	//eventQ is a go channel that
	//into which one goroutine can
	//write byte slices, and out of which
	//another goroutine can read those byte slices
	// eventQ chan []byte

	//TODO: add other fields to this struct
	//that you might need. For example, you'll
	//need to track all of the current WebSocket
	//connections. Remember that slice will be used
	//by multiple goroutines, so you'll need to
	//protect it for concurrent use!

	// map for updating the task list
	connections map[int64]*websocket.Conn
	mx          sync.RWMutex
}

//NewNotifier constructs a new Notifier
func NewNotifier() *Notifier {
	n := &Notifier{
		// eventQ: make(chan []byte, 1024), //buffered channel that can hold 1024 slices at a time
		connections: make(map[int64]*websocket.Conn),
		mx:          sync.RWMutex{},
	}

	//TODO: call the .start() method on its own goroutine
	// go n.start()
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

// //Notify broadcasts the event to all WebSocket clients
// func (n *Notifier) Notify(event []byte) {
// 	log.Printf("adding event to the queue")
// 	//TODO: add `event` to the `n.eventQ`
// 	//see https://tour.golang.org/concurrency/2
// 	//and https://gobyexample.com/channels
// 	n.eventQ <- event
// }

// var buffer = {
//     "type": "",
//     "task": {},
//     "tasks": [],
//     "point": 0,
//     "User": {},
// };

type task struct {
}

type message struct {
	name  string
	task  *task
	tasks []*task
	point int
	user  users.User
}

//start starts the notification loop
func (n *Notifier) Start(msgs <-chan amqp.Delivery) {
	log.Println("starting notifier loop")
	// for msg := range msgs {
		n.mx.Lock()
		defer n.mx.Unlock()
		// for conn := range n.connections {
		// 	var message message
		// 	err := json.Unmarshal(msg.Body, message)
		// 	if err != nil {
		// 		//throw error
		// 	}
		// 	//proabably another forloop?
		// 	// if message.id == conn {
		// 	// 	//write message
		// 	// 	//if there is an error while writing message remove connection
		// 	// }
		// }
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
