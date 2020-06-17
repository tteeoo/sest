package handler

import (
	"fmt"
	"net/http"
)

// AuthHandler handles the /auth endpoint
func AuthHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "{\"message\": \"authentication required\"}")
}
