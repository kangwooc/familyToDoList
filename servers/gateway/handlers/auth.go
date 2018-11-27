package handlers

import (
	"encoding/json"
	"fmt"
	"homework-juan3674-1532739/servers/gateway/models/users"
	"homework-juan3674-1532739/servers/gateway/sessions"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

//UsersHandler handles requests for the "users" resource.
//it will accept POST requests to create new user account
func (context *HandlerContext) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		header := r.Header.Get("Content-Type")
		if !strings.HasPrefix(header, "application/json") {
			http.Error(w, "Request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		var newUser users.NewUser
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			http.Error(w, "Decoding problem", http.StatusBadRequest)
			return
		}

		user, err := newUser.ToUser()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dbuser, err := context.User.Insert(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//Respond to the client
		newSessionState := &SessionState{
			SessionBegan: time.Now(),
			User:         dbuser,
		}
		_, err = sessions.BeginSession(context.SigningKey, context.Session, newSessionState, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Current status method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

//SpecificUserHandler handles requests for a specific user.
func (context *HandlerContext) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id := path.Base(r.URL.Path)
		split := strings.Split(r.URL.Path, "/")
		if len(split) > 4 {
			http.Error(w, "User must be authenticated", http.StatusUnauthorized)
			return
		}
		sessionState := &SessionState{}
		_, err := sessions.GetState(r, context.SigningKey, context.Session, sessionState)
		if err != nil {
			http.Error(w, "User must be authenticated", http.StatusUnauthorized)
			return
		}
		var numID int64
		user := &users.User{}
		if id != "me" {
			numID, err = strconv.ParseInt(id, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			numID = sessionState.User.ID
		}
		user, err = context.User.GetByID(numID)
		if err != nil {
			http.Error(w, "This user is not found in the store", http.StatusNotFound)
			return
		}
		//responding to the client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Current status method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

//SessionHandler handles requests for the "sessions" resource, and allows clients to begin a new session
//using existing user's credentials
func (context *HandlerContext) SessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		header := r.Header.Get("Content-Type")
		if !strings.HasPrefix(header, "application/json") {
			http.Error(w, "Request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		user := &users.Credentials{}
		err := json.NewDecoder(r.Body).Decode(user)
		if err != nil {
			http.Error(w, "Request body could not be parsed", http.StatusInternalServerError)
			return
		}

		profile, err := context.User.GetByEmail(user.Email)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid credentials: %v", err), http.StatusUnauthorized)
			return
		}
		err = profile.Authenticate(user.Password)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid credentials: %v", err), http.StatusUnauthorized)
			return
		}

		newSessionState := &SessionState{
			SessionBegan: time.Now(),
			User:         profile,
		}
		_, err = sessions.BeginSession(context.SigningKey, context.Session, newSessionState, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(profile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Current status method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

//SpecificSessionHandler handles requestss related to a specific authenticated session
func (context *HandlerContext) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		currentSession := path.Base(r.URL.Path)
		if currentSession != "mine" {
			http.Error(w, "Not valid user session", http.StatusForbidden)
			return
		}
		_, err := sessions.EndSession(r, context.SigningKey, context.Session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Signed Out"))
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Current status method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}
