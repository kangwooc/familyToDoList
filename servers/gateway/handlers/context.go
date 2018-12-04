package handlers

import (
	"final-project-zco/servers/gateway/models/users"
	"final-project-zco/servers/gateway/sessions"
)

//HandlerContext is a receiver on any of the HTTP
//handler functions that need access too
//gloabals, such as the key used for signing and veryfying
//SessionIDs, the session store and the user store.
type HandlerContext struct {
	SigningKey string
	Session    sessions.Store
	User       users.Store
	Family     users.Store
	Notifier   Notifier
	Request    map[int64][]*users.User
}
