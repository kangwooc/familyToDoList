package handlers

import (
	"encoding/json"
	"final-project-zco/servers/gateway/sessions"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
)

// NewServiceProxy returns a new ReverseProxy
// for a microservice given a comma-delimited
// list of network addresses
func (ctx *HandlerContext) NewServiceProxy(addrs string) *httputil.ReverseProxy {
	splitAddrs := strings.Split(addrs, ",")
	nextAddr := 0
	mx := sync.Mutex{}

	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			mx.Lock()
			r.URL.Host = splitAddrs[nextAddr]
			nextAddr = (nextAddr + 1) % len(splitAddrs)
			mx.Unlock()

			ss := &SessionState{}
			if _, err := sessions.GetState(r, ctx.SigningKey, ctx.Session, ss); err != nil {
				log.Printf(fmt.Sprintf("session id error: %v", err.Error()))
				return
			}
			userJSON, err := json.Marshal(ss.User)
			if err != nil {
				log.Printf(fmt.Sprintf("json marshal error: %v", err.Error()))
				return
			}
			log.Printf("this is user json %v", userJSON)
			r.Header.Del("X-User")
			r.Header.Set("X-User", string(userJSON))
		},
	}
}
