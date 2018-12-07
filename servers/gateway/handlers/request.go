package handlers

import (
	"bytes"
	"encoding/json"
	"final-project-zco/servers/gateway/models/users"
	"final-project-zco/servers/gateway/sessions"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

type status struct {
	Role     string `json:"personrole"`
	RoomName string `json:"roomname"`
	MemberID int64  `json:"memberid"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// JoinHandler join a family room
// post /join
func (context *HandlerContext) JoinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// header := r.Header.Get("Content-Type")
		// if !strings.HasPrefix(header, "application/json") {
		// 	http.Error(w, "Request body must be in JSON", http.StatusUnsupportedMediaType)
		// 	return
		// }
		sessionState := &SessionState{}
		sid, err := sessions.GetState(r, context.SigningKey, context.Session, sessionState)
		if err != nil {
			http.Error(w, "User must be authenticated", http.StatusUnauthorized)
			return
		}
		numID := sessionState.User.ID
		// if admin, no allowed to join
		if sessionState.User.Role == "Admin" {
			http.Error(w, "Admin can not join other room", http.StatusBadRequest)
			return
		}
		var update *users.Updates
		// decode the entered family room name
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		admin, err := context.User.GetAdmin(update.RoomName, "Admin")
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		member := &users.Updates{RoomName: update.RoomName, Role: "Waiting"}
		// update the user role to be admin
		added, err := context.User.Update(numID, member)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sessionState.User.Role = member.Role
		sessionState.User.RoomName = member.RoomName
		if err = context.Session.Save(sid, sessionState); err != nil {

			log.Printf(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = added.ApplyUpdates(member); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newSlice, ok := context.Request[admin.ID]
		var userSlice []*users.User
		if !ok {
			userSlice = make([]*users.User, 0)
		} else {
			userSlice = newSlice
		}

		rabbit := os.Getenv("RABBITADDR")

		conn, err := amqp.Dial("amqp://" + rabbit + "/")
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()
		q, err := ch.QueueDeclare(
			"authQueue", // name
			false,       // durable
			false,       // delete when unused
			false,       // exclusive
			false,       // no-wait
			nil,         // arguments
		)
		failOnError(err, "Failed to declare a queue")

		var user bytes.Buffer // Stand-in for a network connection
		body := sessionState.User
		err = json.NewEncoder(&user).Encode(body)
		if err != nil {
			log.Fatal("encode error:", err)
		}
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        user.Bytes(),
			})
		failOnError(err, "Failed to publish a message")

		context.Request[admin.ID] = append(userSlice, added)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Sent"))
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Current status method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// ReceiveHandler is the method that receive the request
func (context *HandlerContext) ReceiveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		//Check authority and get context.Request if it's empty return empty json.
		header := r.Header.Get("Content-Type")
		if !strings.HasPrefix(header, "application/json") {
			http.Error(w, "Request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		sessionState := &SessionState{}
		_, err := sessions.GetState(r, context.SigningKey, context.Session, sessionState)
		if err != nil {
			http.Error(w, "User must be authenticated", http.StatusUnauthorized)
			return
		}
		switch sessionState.User.Role {
		case "Admin":
			break
		case "Waiting":
			break
		default:
			http.Error(w, "Admin and Waiting can get", http.StatusUnauthorized)
			return
			break

		}
		numID := sessionState.User.ID
		request, ok := context.Request[numID]
		var result []*users.User
		if !ok {
			result = make([]*users.User, 0)
		} else {
			result = request
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Current status method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// AcceptRequest is the method that admin can accept the requests
func (context *HandlerContext) AcceptRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		header := r.Header.Get("Content-Type")
		if !strings.HasPrefix(header, "application/json") {
			http.Error(w, "Request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		sessionState := &SessionState{}
		if _, err := sessions.GetState(r, context.SigningKey, context.Session, sessionState); err != nil {
			http.Error(w, "User must be authenticated", http.StatusUnauthorized)
			return
		}
		if sessionState.User.Role != "Admin" {
			http.Error(w, "Admin only can get", http.StatusUnauthorized)
			return
		}
		var accept status
		if err := json.NewDecoder(r.Body).Decode(&accept); err != nil {
			http.Error(w, "Decoding problem", http.StatusBadRequest)
			return
		}
		up := &users.Updates{Role: accept.Role, RoomName: accept.RoomName}
		added, err := context.User.UpdateToMember(accept.MemberID, up)
		log.Printf("Debug: Admin: ", sessionState.User)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = added.ApplyUpdates(up); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//get slice of user from the map
		//remove it from the slice.
		var slice []*users.User
		slice = make([]*users.User, 0)
		if len(context.Request[sessionState.User.ID]) == 1 {
			context.Request[sessionState.User.ID] = make([]*users.User, 0)
		} else {
			for _, k := range context.Request[sessionState.User.ID] {
				if k.ID != accept.MemberID {
					//append
					slice = append(slice, k)
				}
			}
			context.Request[sessionState.User.ID] = slice
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Request complete!"))
	} else {
		http.Error(w, "Current status method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}
