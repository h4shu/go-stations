package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

func Logging(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			t := time.Now()
			elapsed := t.Sub(start)
			ctx := r.Context()
			o, ok := ctx.Value(model.ContextKey("OS")).(string)
			if !ok {
				log.Println("ContextKey OS not found")
				return
			}
			l := &model.AccessLog{
				Timestamp: start,
				Latency:   elapsed.Milliseconds(),
				Path:      r.URL.Path,
				OS:        o,
			}
			b, err := json.Marshal(l)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(string(b))
			// JSON に変換せず以下で出力可
			// fmt.Printf("%+v\n", l)
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
