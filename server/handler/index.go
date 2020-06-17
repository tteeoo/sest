package handler

import (
	"fmt"
	"net/http"
)

// IndexHandler handles the / page
func IndexHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		ErrorHandler(w, r, http.StatusNotFound)
		return
	}

	fmt.Fprint(w, "found")
}
