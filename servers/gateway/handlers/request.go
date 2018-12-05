package handlers

import (
	"encoding/json"
	"final-project-zco/servers/gateway/models/users"
	"final-project-zco/servers/gateway/sessions"
	"log"
	"net/http"
	"strings"
)

type status struct {
	Role     string `json:"personrole"`
	RoomName string `json:"roomname"`
	MemberID int64  `json:"memberid"`
}

// JoinHandler join a family room
// post
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

		member := &users.Updates{RoomName: update.RoomName, Role: "Waiting"}
		log.Printf("this is qqq id %d", numID)
		log.Printf("this is roomname from userupdate %d", update.RoomName)
		added, err := context.User.Update(numID, member)
		if err != nil {
			log.Printf("aaa")

			log.Printf(err.Error())

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sessionState.User.Role = member.Role
		sessionState.User.RoomName = member.RoomName
		if err = context.Session.Save(sid, sessionState); err != nil {
			log.Printf("adddaa")

			log.Printf(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = added.ApplyUpdates(member); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		admin, err := context.User.GetAdmin(update.RoomName, "Admin")
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		newSlice, ok := context.Request[admin.ID]
		var userSlice []*users.User
		if !ok {
			userSlice = make([]*users.User, 0)
		} else {
			userSlice = newSlice
		}
		context.Request[admin.ID] = append(userSlice, added)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Sent"))
		//Is this right status?
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Current status method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// get
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
		if sessionState.User.Role != "Admin" {
			http.Error(w, "Admin only can get", http.StatusUnauthorized)
			return
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

// post
func (context *HandlerContext) AcceptRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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
		if sessionState.User.Role != "Admin" {
			http.Error(w, "Admin only can get", http.StatusUnauthorized)
			return
		}
		var accept status
		log.Printf("this is body yo %v", r.Body)

		err = json.NewDecoder(r.Body).Decode(&accept)
		if err != nil {
			http.Error(w, "Decoding problem", http.StatusBadRequest)
			return
		}
		up := &users.Updates{Role: "Member", RoomName: accept.RoomName}
		added, err := context.User.UpdateToMember(accept.MemberID, up)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// sessionState.added.Role = "Member"
		// sessionState.added.RoomName = accept.RoomName
		// if err = context.Session.Save(sid, sessionState); err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		if err = added.ApplyUpdates(up); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		q, _ := context.User.GetByID(accept.MemberID)

		log.Printf("mem id %v", q)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Request complete!"))
	} else {
		http.Error(w, "Current status method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}
