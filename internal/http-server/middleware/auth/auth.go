package mwAuth

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	resp "github.com/solumD/go-blog-api/internal/lib/api/response"
	"github.com/solumD/go-blog-api/internal/lib/jwt"
	"github.com/solumD/go-blog-api/internal/lib/logger/sl"
)

type Response struct {
	resp.Response
}

// написать логику для проверки jwt токена
func New(secret string, log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/auth"),
		)

		fn := func(w http.ResponseWriter, r *http.Request) {
			log.Info("middleware auth used")

			auth := r.Header.Get("Authorization")
			token := strings.Split(auth, " ")[1]

			claims, err := jwt.GetTokenClaims(secret, token)
			if err != nil {
				log.Error("invalid jwt token", sl.Err(err))

				render.JSON(w, r, resp.Error(fmt.Sprintf("%v", err)))

				return
			}

			login := claims["sub"].(string)
			r.Header.Add("Login", login)

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
