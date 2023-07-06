package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO: ここに実装をする
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}

			res := &model.HealthzResponse{Message: "Panicking"}
			err := json.NewEncoder(w).Encode(res)
			if err != nil {
				log.Println(err)
			}
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
