package handlers

import (
	"encoding/json"
	"final-project-zco/servers/gateway/models/users"
	"final-project-zco/servers/gateway/sessions"
	"net/http"
	"strings"
)

type status struct {
	User     *users.Updates
	MemberID int64
}

// JoinHandler join a family room
func (context *HandlerContext) JoinHandler(w http.ResponseWriter, r *http.Request) {
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
		numID := sessionState.User.ID

		var update *users.Updates
		// decode the entered family room name
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		member := &users.Updates{RoomName: update.RoomName, Role: update.Role}
		// update the user role to be admin
		if _, err := context.User.Update(numID, member); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		admin, err := context.User.GetAdmin(update.RoomName, "admin")
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
		context.Request[admin.ID] = append(userSlice, sessionState.User)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Sent"))
		//Is this right status?
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Current status method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (context *HandlerContext) receiveHandler(w http.ResponseWriter, r *http.Request) {
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
		err = json.NewDecoder(r.Body).Decode(accept)
		if err != nil {
			http.Error(w, "Decoding problem", http.StatusBadRequest)
			return
		}

		if _, err := context.User.UpdateToMember(accept.MemberID, accept.User); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Current status method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}
