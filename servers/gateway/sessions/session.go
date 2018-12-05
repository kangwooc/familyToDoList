package sessions

import (
	"errors"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	id, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	// log.Printf("begin session: %v", sessionState)
	if err = store.Save(id, sessionState); err != nil {
		return InvalidSessionID, err
	}
	w.Header().Add(headerAuthorization, schemeBearer+string(id))
	return id, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	// Reference: https://gist.github.com/sambengtson/bc9f76331065f09e953f
	var str string
	header := strings.SplitN(r.Header.Get(headerAuthorization), " ", 2)
	if len(header) == 0 {
		// If no Authorization header is present, find the "auth" query string parameter and split into 2
		authQuery := strings.SplitN(r.FormValue(paramAuthorization), " ", 2)
		// If it doesn't have anything on neither, return empty SessionID and error called ErrNoSessionID
		if len(authQuery) < 2 {
			return InvalidSessionID, ErrNoSessionID
		}
		str = authQuery[1]
	} else {
		// Split the Authorization hearer into 2 (e.g. ["Bearer", "SessionID"])
		if header[0] != "Bearer" {
			return InvalidSessionID, ErrInvalidScheme
		}
		str = header[1]
	}
	// Validate sessionID from Authorization header
	sid, err := ValidateID(str, signingKey)
	// If it's invalid scheme, return sid and error called ErrInvalidScheme
	if err != nil {
		return InvalidSessionID, err
	}
	// If it's valid, return the SessionID.
	return sid, nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	id, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, ErrStateNotFound
	}
	if err = store.Get(id, sessionState); err != nil {
		return InvalidSessionID, ErrStateNotFound
	}
	// log.Printf("get state: %v", sessionState)

	return id, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	id, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, err
	}
	if err = store.Delete(SessionID(id)); err != nil {
		return InvalidSessionID, err
	}
	return SessionID(id), nil
}
