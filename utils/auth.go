package utils

import (
	"net/http"
)

func GetUser(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("access-token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Cookie not found", http.StatusUnauthorized)
			return ""
		}
		http.Error(w, "", http.StatusInternalServerError)
		return ""
	}

	return cookie.Value
}
