package middleware

import (
	"net/http"
	"encoding/json"
	"time"
	"log"
	"fmt"

	"github.com/TechBowl-japan/go-stations/model"
)

func Logging(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			t := time.Now()
			elapsed := t.Sub(start)
			c := r.Context()
			o, ok := c.Value(model.ContextKey("OS")).(string)
			if !ok {
				log.Println("ContextKey OS not found")
			} else {
				l := &model.AccessLog{
					Timestamp: start,
					Latency: elapsed.Milliseconds(),
					Path: r.URL.Path,
					OS: o,
				}
				b, err := json.Marshal(l)
				if err != nil {
					log.Println(err)
				}
				fmt.Println(string(b))
			}
		}()

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}