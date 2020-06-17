package handler

import (
	"fmt"
	"github.com/tteeoo/go-website/util"
	"net/http"
	"strconv"
)

// ErrorHandler handles errors by taking a status code and rendering a template with text
func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {

	util.Logger.Println("ERROR: " + util.GetRemoteAddr(r) + " " + strconv.Itoa(status))
	w.WriteHeader(status)
	fmt.Fprint(w, status, " " + http.StatusText(status))
}
