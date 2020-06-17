package handler

import (
	"fmt"
	"github.com/tteeoo/sest/server/util"
	"net/http"
	"strconv"
)

// ErrorHandler handles errors by taking a status code and rendering a template with text
func ErrorHandler(w http.ResponseWriter, r *http.Request, status int) {

	util.Logger.Println("sest-server: http-error: host: " + util.GetRemoteAddr(r) + " error: " + strconv.Itoa(status) + " uri: " + r.RequestURI + " method: " + r.Method + " ua: " + r.UserAgent())
	w.WriteHeader(status)
	fmt.Fprint(w, "{\"message:\": \""+strconv.Itoa(status)+" "+http.StatusText(status)+"\"}")
}
