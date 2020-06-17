package handler

import (
	"github.com/tteeoo/sest/server/limit"
	"github.com/tteeoo/sest/server/util"
	"net/http"
)

var limiter = limit.NewIPRateLimiter(1, 5)

func GlobalHandler(handle func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Always respond with json
		w.Header().Add("Content-Type", "application/json")

		// Check for rate-limiting
		limiter := limiter.GetLimiter(util.GetRemoteAddr(r))
		if !limiter.Allow() {
			ErrorHandler(w, r, http.StatusTooManyRequests)
			return
		}

		// Check black/whitelist
		if !util.CheckAllowed(util.GetRemoteAddr(r), Whitelist, Blacklist) {
			ErrorHandler(w, r, http.StatusForbidden)
			return
		}

		// Logging
		util.Logger.Println("sest-server: request: host: " + util.GetRemoteAddr(r) + " uri: " + r.RequestURI + " method: " + r.Method + " ua: " + r.UserAgent())

		handle(w, r)
	}
}
