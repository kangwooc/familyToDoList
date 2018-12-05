package handlers

import (
	"encoding/json"
	"final-project-zco/servers/gateway/models/users"
	"final-project-zco/servers/gateway/sessions"
	"log"
	"net/http"
	"path"
	"strconv"
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
		user := &users.User{}
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		update := &users.Updates{Role: "Default", RoomName: ""}
		u, err := context.User.UpdateToMember(user.ID, update)
		if err != nil {
			// log.Printf("Debug: sessionState User Role on Delete: ", sessionState.User.Role)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = u.ApplyUpdates(update); err != nil {
			// log.Printf("what is wrofffng222", sessionState.User.Role)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
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
			log.Printf("rrrrole %v", sessionState.User.Role)
			http.Error(w, "User must be admin to view all member", http.StatusUnauthorized)
			return
		}
		// if user is authenticated, get the room id
		roomid := path.Base(r.URL.Path)
		log.Println("this is roomid %s", roomid)
		room, err := strconv.ParseInt(roomid, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		split := strings.Split(r.URL.Path, "/")
		log.Println("this is split %d", len(split))
		if len(split) > 3 {
			http.Error(w, "User must be authenticated", http.StatusUnauthorized)
			return
		}
		fam, err := context.Family.GetRoomName(room)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("this is fam %v", fam.RoomName)

		// once get the room name, get all the users in that room
		userArr, err := context.User.GetByRoomName(fam.RoomName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("this is userarr %v", userArr)

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
