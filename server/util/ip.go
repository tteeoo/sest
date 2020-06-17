package util

import (
	"net/http"
	"strings"
)

// GetRemoteAddr will get the remote address taking into account Cloudflare
func GetRemoteAddr(r *http.Request) string {

	ip := r.Header.Get("CF-Connecting-IP")
	if ip == "" {
		return r.RemoteAddr
	}

	return ip
}

// CheckAllowed will check a host against the black/whitelist and return true if they are allowed access
func CheckAllowed(host string, whitelist, blacklist []string) bool {

	if len(whitelist) > 0 {
		for _, i := range whitelist {
			var colons int
			for _, j := range i {
				if j == ':' {
					colons++
				}
			}
			if colons < 2 {
				if strings.Split(host, ":")[0] == strings.Split(i, ":")[0] {
					return true
				}
			} else {
				if host == i {
					return true
				}
			}
		}
		return false
	} else if len(blacklist) > 0 {
		for _, i := range blacklist {
			var colons int
			for _, j := range i {
				if j == ':' {
					colons++
				}
			}
			if colons < 2 {
				if strings.Split(host, ":")[0] == strings.Split(i, ":")[0] {
					return false
				}
			} else {
				if host == i {
					return false
				}
			}
		}
		return true
	}

	return true
}
