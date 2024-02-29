package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"playlistturbo.com/model"
	"playlistturbo.com/plterror"
)

var excludeFromLogger = []string{
	"/api/notes/checkupdates",
}

// LoggerMiddleware log success and error responses details.
func LoggerMiddleware() Middleware {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			user := "PUBLIC" // default
			start := time.Now()

			id, ok := r.Context().Value(&model.CtxKeyID).(uuid.UUID)
			if !ok {
				log.Println("panic LoggerMiddleware")
				panic(plterror.ErrServerError)
			}

			isPathOk := true
			for _, path := range excludeFromLogger {
				if path == r.URL.Path {
					isPathOk = false
				}
			}

			if isPathOk {
				defer func() {
					elapsedTime := time.Since(start).String()
					plterror.Logger.Infof("%s %s %s %s %s", id.String(), r.Method, r.URL.Path, user, elapsedTime)
				}()
			}

			h.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	return f
}
