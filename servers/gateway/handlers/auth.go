package handlers

import (
	"encoding/json"
	"final-project-zco/servers/gateway/models/users"
	"final-project-zco/servers/gateway/sessions"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

//UsersHandler handles requests for the "users" resource.
//it will accept POST requests to create new user account
// sign up => /users
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

// CreateHandler create a family room
// post "/create"
func (context *HandlerContext) CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		header := r.Header.Get("Content-Type")
		if !strings.HasPrefix(header, "application/json") {
			http.Error(w, "Request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		// id := path.Base(r.URL.Path)
		// split := strings.Split(r.URL.Path, "/")
		// if len(split) > 4 {
		// 	http.Error(w, "User must be authenticated", http.StatusUnauthorized)
		// 	return
		// }
		sessionState := &SessionState{}
		sid, err := sessions.GetState(r, context.SigningKey, context.Session, sessionState)
		if err != nil {
			http.Error(w, "User must be authenticated", http.StatusUnauthorized)
			return
		}

		numID := sessionState.User.ID
		var family *users.FamilyRoom
		if err := json.NewDecoder(r.Body).Decode(&family); err != nil {
			http.Error(w, "Decoding problem", http.StatusBadRequest)
			return
		}
		admin := &users.Updates{Role: "Admin", RoomName: family.RoomName}
		// update the user role to be admin
		added, err := context.User.UpdateToMember(numID, admin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sessionState.User.Role = admin.Role
		sessionState.User.RoomName = admin.RoomName
		if err = context.Session.Save(sid, sessionState); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err = added.ApplyUpdates(admin); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("this is user yo %v", added)
		// insert into family table
		fam, err := context.Family.InsertFam(family)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err = json.NewEncoder(w).Encode(fam); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Current status method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// // JoinHandler join a family room
// func (context *HandlerContext) JoinHandler(w http.ResponseWriter, r *http.Request) {
// 	log.Println("heihhh")
// 	if r.Method == http.MethodPatch { // what method
// 		header := r.Header.Get("Content-Type")
// 		if !strings.HasPrefix(header, "application/json") {
// 			http.Error(w, "Request body must be in JSON", http.StatusUnsupportedMediaType)
// 			return
// 		}
// 		// id := path.Base(r.URL.Path)
// 		// split := strings.Split(r.URL.Path, "/")
// 		// if len(split) > 4 {
// 		// 	http.Error(w, "User must be authenticated", http.StatusUnauthorized)
// 		// 	return
// 		// }
// 		sessionState := &SessionState{}
// 		_, err := sessions.GetState(r, context.SigningKey, context.Session, sessionState)
// 		if err != nil {
// 			http.Error(w, "User must be authenticated", http.StatusUnauthorized)
// 			return
// 		}
// 		numID := sessionState.User.ID

// 		var update *users.Updates
// 		// decode the entered family room name
// 		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		member := &users.Updates{RoomName: update.RoomName}
// 		// update the user role to be admin
// 		if _, err := context.User.UpdateToMember(numID, member); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 	} else {
// 		http.Error(w, "Current status method is not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}
// }

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
		if id != "me" {
			numID, err = strconv.ParseInt(id, 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			numID = sessionState.User.ID
		}
		user, err := context.User.GetByID(numID)
		if err != nil {
			http.Error(w, "This user do not found in the store", http.StatusNotFound)
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
// sign in
func (context *HandlerContext) SessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		header := r.Header.Get("Content-Type")
		if !strings.HasPrefix(header, "application/json") {
			http.Error(w, "Request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}
		user := &users.Credentials{}
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			http.Error(w, "Request body could not be parsed", http.StatusInternalServerError)
			return
		}

		profile, err := context.User.GetByUserName(user.UserName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid credentials: %v", err), http.StatusUnauthorized)
			return
		}
		if err = profile.Authenticate(user.Password); err != nil {
			http.Error(w, fmt.Sprintf("Invalid credentials: %v", err), http.StatusUnauthorized)
			return
		}

		newSessionState := &SessionState{
			SessionBegan: time.Now(),
			User:         profile,
		}
		if _, err = sessions.BeginSession(context.SigningKey, context.Session, newSessionState, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err = json.NewEncoder(w).Encode(profile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Current status method is not allowed", http.StatusMethodNotAllowed)
		return
	}
}

//SpecificSessionHandler handles requestss related to a specific authenticated session
// sign out
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
