package middleware

import (
	"net/http"
	"os"
)

func BasicAuth(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		uid, pass, ok := r.BasicAuth()
		if !ok || uid != os.Getenv("BASIC_AUTH_USER_ID") || pass != os.Getenv("BASIC_AUTH_PASSWORD") {
			w.Header().Add("WWW-Authenticate", `Basic realm="SECRET AREA"`)
			http.Error(w, "Not authorized.", 401)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
