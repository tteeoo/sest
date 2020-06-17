package handler

import (
	"fmt"
	"net/http"
)

// IndexHandler handles the / page
func IndexHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "found")
}
