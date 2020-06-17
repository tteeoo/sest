package main

import (
	"github.com/tteeoo/sest/server/handler"
	"github.com/tteeoo/sest/server/limit"
	"github.com/tteeoo/sest/server/util"
	"net/http"
	"os"
)

var cont string
var addr = os.Getenv("SEST_SERVER_ADDR")
var contDir = os.Getenv("SEST_SERVER_CONT_DIR")

var limiter = limit.NewIPRateLimiter(1, 5)

func init() {

	// Default env vars
	if len(addr) == 0 {
		addr = "127.0.0.1:7000"
	}
	if len(contDir) == 0 {
		contDir = os.Getenv("HOME") + "/.sest"
	}

	// Get continer
	if len(os.Args) < 2 {
		util.Logger.Fatal("sest-server: error: no container specified")
	}
	cont = os.Args[1]
}

func main() {

	// Make dir if it does not exist
	if _, err := os.Stat(contDir); os.IsNotExist(err) {
		util.Logger.Println("sest-server: making directory at ", contDir)
		os.Mkdir(contDir, 0700)
	}

	// Setup logger
	defer util.LogFile.Close()

	// Handle routes
	http.HandleFunc("/", globalHandler(handler.IndexHandler))

	// Start the server
	util.Logger.Println("sest-server: attempting to listen on http://" + addr)
	util.Logger.Fatal("sest-server: error: ", http.ListenAndServe(addr, nil))
}

func globalHandler(handle func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Always respond with json
		w.Header().Add("Content-Type", "application/json")

		// Check for rate-limiting
		limiter := limiter.GetLimiter(util.GetRemoteAddr(r))
		if !limiter.Allow() {
			handler.ErrorHandler(w, r, http.StatusTooManyRequests)
			return
		}

		// Logging
		util.Logger.Println("sest-server: request: host: " + util.GetRemoteAddr(r) + " uri: " + r.RequestURI + " method: " + r.Method + " ua: " + r.UserAgent())

		handle(w, r)
	}
}
