package middleware

import (
	"context"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/mileusna/useragent"
)

func UA(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ua := useragent.Parse(r.UserAgent())
		ctx := r.Context()
		*r = *r.WithContext(context.WithValue(ctx, model.ContextKey("OS"), ua.OS))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
