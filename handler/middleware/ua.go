package middleware

import (
	"net/http"
	"context"

	"github.com/mileusna/useragent"
	"github.com/TechBowl-japan/go-stations/model"
)

func UA(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ua := useragent.Parse(r.UserAgent())
		c := r.Context()
		*r = *r.WithContext(context.WithValue(c, model.ContextKey("OS"), ua.OS))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}