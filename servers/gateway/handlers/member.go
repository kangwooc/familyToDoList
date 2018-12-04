package handlers

import (
	"encoding/json"
	"final-project-zco/servers/gateway/models/users"
	"final-project-zco/servers/gateway/sessions"
	"log"
	"net/http"
	"strings"
)

//fn ln id photourl

// delete member
func (context *HandlerContext) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		// get the member info
		header := r.Header.Get("Content-Type")
		if !strings.HasPrefix(header, "application/json") {
			http.Error(w, "Request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}

		sessionState := &SessionState{}
		sid, err := sessions.GetState(r, context.SigningKey, context.Session, sessionState)
		if err != nil {
			http.Error(w, "User must be authenticated", http.StatusUnauthorized)
			return
		}

		// check whether current user is an admin
		if sessionState.User.Role != "Admin" {
			log.Printf("rrrrole %v", sessionState.User.Role)
			http.Error(w, "User must be admin to delete member", http.StatusUnauthorized)
			return
		}
		// the member to delete
		user := &users.User{}
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		update := &users.Updates{Role: "Default", RoomName: ""}
		u, err := context.User.UpdateToMember(user.ID, update)
		if err != nil {
			log.Printf("what is wrong", sessionState.User.Role)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sessionState.User.Role = update.Role
		sessionState.User.RoomName = update.RoomName
		if err = context.Session.Save(sid, sessionState); err != nil {
			log.Printf("what is wrong222", sessionState.User.Role)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = u.ApplyUpdates(update)
		if err != nil {
			log.Printf("what is wrofffng222", sessionState.User.Role)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(u); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Current status method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// display all members for each room
