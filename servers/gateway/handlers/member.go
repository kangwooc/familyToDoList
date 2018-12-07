package handlers

import (
	"encoding/json"
	"final-project-zco/servers/gateway/models/users"
	"final-project-zco/servers/gateway/sessions"
	"log"
	"net/http"
	"path"
	"strings"
)

type Num struct {
	ID int `json:"id"`
}

// DeleteHandler is for deleting member
func (context *HandlerContext) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		// get the member info
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

		// check whether current user is an admin
		if sessionState.User.Role != "Admin" {
			http.Error(w, "User must be admin to delete member", http.StatusUnauthorized)
			return
		}
		// the member to delete
		num := &Num{}
		if err := json.NewDecoder(r.Body).Decode(num); err != nil {
			log.Printf("resulttthehe %v", num)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		update := &users.Updates{Role: "Default", RoomName: ""}
		u, err := context.User.UpdateToMember(int64(num.ID), update)
		if err != nil {
			log.Printf("this is id last %v", num.ID)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = u.ApplyUpdates(update); err != nil {
			log.Printf("what is wrofffng222", sessionState.User.Role)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Deleted"))
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Current status method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// DisplayHandler is displaying all members for each room
// get localhost/room/1
func (context *HandlerContext) DisplayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// get the member info
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
		// check whether current user is an admin
		if sessionState.User.Role != "Admin" {
			// log.Printf("rrrrole %v", sessionState.User.Role)
			http.Error(w, "User must be admin to view all member", http.StatusUnauthorized)
			return
		}
		// if user is authenticated, get the room id
		roomname := path.Base(r.URL.Path)
		// once get the room name, get all the users in that room
		userArr, err := context.User.GetByRoomName(roomname)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(userArr); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Current status method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}
