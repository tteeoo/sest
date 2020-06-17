package util

import (
	"net/http"
)

// GetRemoteAddr will get the remote address taking into account Cloudflare
func GetRemoteAddr(r *http.Request) string {

	ip := r.Header.Get("CF-Connecting-IP")
	if ip == "" {
		return r.RemoteAddr
	}

	return ip
}
