package mwAuth

import (
	"log/slog"
	"net/http"
)

// написать логику для проверки jwt токена
func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/auth"),
		)

		fn := func(w http.ResponseWriter, r *http.Request) {
			log.Info("middleware auth used")

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
