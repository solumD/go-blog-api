package login

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/solumD/go-blog-api/internal/lib/api/response"
	"github.com/solumD/go-blog-api/internal/lib/jwt"
	"github.com/solumD/go-blog-api/internal/lib/logger/sl"
	"github.com/solumD/go-blog-api/internal/lib/password"
)

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	resp.Response
	Token string `json:"token,omitempty"`
}

type UserAuthorizer interface {
	IsUserExist(login string) (bool, error)
	GetPassword(login string) (string, error)
}

func New(secret string, log *slog.Logger, userAuthorizer UserAuthorizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.login.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		req.Login = strings.TrimSpace(req.Login)
		req.Password = strings.TrimSpace(req.Password)

		exist, err := userAuthorizer.IsUserExist(req.Login)
		if err != nil {
			log.Error("failed to check if user exists", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to check if user exists"))

			return
		}

		if !exist {
			log.Error("invalid request", sl.Err(fmt.Errorf("user does not exist")))

			render.JSON(w, r, resp.Error("user does not exist"))

			return
		}

		realPassword, err := userAuthorizer.GetPassword(req.Login)
		if err != nil {
			log.Error("failed to get user's real password", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to get user's real password"))

			return
		}

		if err := password.CompareHashAndPass(req.Password, realPassword); err != nil {
			log.Error("invalid request", sl.Err(fmt.Errorf("invalid password")))

			render.JSON(w, r, resp.Error("invalid password"))

			return
		}

		token, err := jwt.GenerateToken(req.Login, secret)
		if err != nil {
			log.Error("failed to generate jwt-token", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to generate jwt-token"))

			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Token:    token,
		})
	}
}
