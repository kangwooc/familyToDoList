package handlers

import (
	"final-project-zco/servers/gateway/models/users"
	"time"
)

//TODO: define a session state struct for this web server
//see the assignment description for the fields you should include
//remember that other packages can only see exported fields!

//SessionState tracks the time at which this session began
//and the authenticated users.User who started the session
type SessionState struct {
	SessionBegan time.Time
	User         *users.User `json:"user"`
}
