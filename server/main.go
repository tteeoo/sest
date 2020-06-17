package main

import (
	"github.com/tteeoo/sest/server/handler"
	"github.com/tteeoo/sest/server/util"
	"net/http"
	"os"
)

var addr = os.Getenv("SEST_SERVER_ADDR")

func init() {

	// Default env vars
	if len(addr) == 0 {
		addr = "127.0.0.1:7000"
	}
}

func main() {

	// Setup logger
	defer util.LogFile.Close()

	// Handle routes
	http.HandleFunc("/", handler.GlobalHandler(func(w http.ResponseWriter, r *http.Request) { handler.ErrorHandler(w, r, http.StatusNotFound) }))
	http.HandleFunc("/auth", handler.GlobalHandler(handler.AuthHandler))

	// Start the server
	util.Logger.Println("sest-server: attempting to listen on http://" + addr)
	util.Logger.Fatal("sest-server: error: ", http.ListenAndServe(addr, nil))
}
