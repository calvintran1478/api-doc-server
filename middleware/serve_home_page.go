package middleware

import (
	"net/http"
)

func ServeHomePage(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			r.URL.Path = "/projects"
		}
		next.ServeHTTP(w, r)
	})
}
