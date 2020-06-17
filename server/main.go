package main

import (
	"fmt"
	"github.com/tteeoo/sest/server/handler"
	"github.com/tteeoo/sest/server/limit"
	"github.com/tteeoo/sest/server/util"
	"net/http"
	"os"
)

var addr string

var limiter = limit.NewIPRateLimiter(1, 5)

func init() {

	// Get addr if set
	addr = os.Getenv("WEB_ADDR")
	if len(addr) == 0 {
		addr = "127.0.0.1:7000"
	}
}

func main() {

	// Setup logger
	defer util.LogFile.Close()

	// Handle routes
	http.HandleFunc("/", rateLimit(handler.IndexHandler))

	// Start the server
	util.Logger.Println("Attempting to listen on http://" + addr)
	util.Logger.Fatal(http.ListenAndServe(addr, nil))
}

func rateLimit(handle func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		limiter := limiter.GetLimiter(util.GetRemoteAddr(r))
		if !limiter.Allow() {
			fmt.Fprint(w, "429 too many requests")
			return
		}

		util.Logger.Println("HIT: " + util.GetRemoteAddr(r) + " " + r.RequestURI)
		handle(w, r)
	}
}
