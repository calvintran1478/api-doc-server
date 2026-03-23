package utils

import (
	"net/http"
)

func GetUser(w http.ResponseWriter, r *http.Request) string {
	cookie, err := r.Cookie("access-token")
	if err != nil {
		http.Redirect(w, r, "/login.html", http.StatusFound)
		return ""
	}

	return cookie.Value
}
