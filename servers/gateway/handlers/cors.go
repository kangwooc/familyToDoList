package handlers

import (
	"net/http"
)

//Cors Struct is used for to build middleware cors handler
type Cors struct {
	MyHandler http.Handler
}

func (c *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Max-Age", "600")
	w.Header().Add(headerAccessControlExposeHeaders, contentAuth)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
	} else {
		c.MyHandler.ServeHTTP(w, r)
	}
}

//NewCors constructs a new cors middleware handler
func NewCors(handlerToWrap http.Handler) *Cors {
	return &Cors{handlerToWrap}
}
