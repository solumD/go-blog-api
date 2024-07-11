package register

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/solumD/go-blog-api/internal/lib/api/response"
	"github.com/solumD/go-blog-api/internal/lib/logger/sl"
	"github.com/solumD/go-blog-api/internal/lib/password"
)

type Request struct {
	Login            string `json:"login"`
	Password         string `json:"password"`
	RepeatedPassword string `json:"repeated_password"`
}

type Response struct {
	resp.Response
	ID int64 `json:"id,omitempty"`
}

type UserRegistrar interface {
	IsUserExist(login string) (bool, error)
	SaveUser(login string, password string) (int64, error)
}

// TODO: добавить валидацию пароля (длина, содержание спецсимволов и тд) через регулярку
func New(log *slog.Logger, userRegistrar UserRegistrar) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.register.New"

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

		log.Info("request body decoded", slog.Any("request", req))
		req.Login = strings.TrimSpace(req.Login)
		req.Password = strings.TrimSpace(req.Password)
		req.RepeatedPassword = strings.TrimSpace(req.RepeatedPassword)

		if fields := strings.Fields(req.Login); len(fields) > 1 {
			log.Error("invalid request", sl.Err(fmt.Errorf("login cannot contain spaces")))

			render.JSON(w, r, resp.Error("login cannot contain spaces"))

			return
		}

		if len(req.Login) < 8 {
			log.Error("invalid request", sl.Err(fmt.Errorf("login cannot be shorter than 8 characters")))

			render.JSON(w, r, resp.Error("login cannot be shorter than 8 characters"))

			return
		}

		exist, err := userRegistrar.IsUserExist(req.Login)
		if err != nil {
			log.Error("failed to check if user exists", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to check if user exists"))

			return
		}

		if exist {
			log.Error("invalid request", sl.Err(fmt.Errorf("user already exists")))

			render.JSON(w, r, resp.Error("user already exists"))

			return
		}

		if len(req.Password) < 8 {
			log.Error("invalid request", sl.Err(fmt.Errorf("password cannot be shorter than 8 characters")))

			render.JSON(w, r, resp.Error("password cannot be shorter than 8 characters"))

			return
		}

		if len(req.RepeatedPassword) < 8 {
			log.Error("invalid request", sl.Err(fmt.Errorf("repeated password cannot be shorter than 8 characters")))

			render.JSON(w, r, resp.Error("repeated password cannot be shorter than 8 characters"))

			return
		}

		if fields := strings.Fields(req.Password); len(fields) > 1 {
			log.Error("invalid request", sl.Err(fmt.Errorf("password cannot contain spaces")))

			render.JSON(w, r, resp.Error("password cannot contain spaces"))

			return
		}

		if fields := strings.Fields(req.RepeatedPassword); len(fields) > 1 {
			log.Error("invalid request", sl.Err(fmt.Errorf("repeated password cannot contain spaces")))

			render.JSON(w, r, resp.Error("repeated password cannot contain spaces"))

			return
		}

		if req.Password != req.RepeatedPassword {
			log.Error("invalid request", sl.Err(fmt.Errorf("password and repeated password don't match")))

			render.JSON(w, r, resp.Error("password and repeated password don't match"))

			return
		}

		hashedPassword, err := password.EncryptPassword(req.Password)
		if err != nil {
			log.Error("failed to encrypt password", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to encrypt password"))

			return
		}

		id, err := userRegistrar.SaveUser(req.Login, hashedPassword)
		if err != nil {
			log.Error("failed to save user", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save user"))

			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
			ID:       id,
		})
	}
}
