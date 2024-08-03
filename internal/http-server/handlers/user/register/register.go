package register

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/solumD/go-blog-api/internal/lib/api/response"
	"github.com/solumD/go-blog-api/internal/lib/logger/sl"
	"github.com/solumD/go-blog-api/internal/lib/password"
	"github.com/solumD/go-blog-api/internal/lib/validator"
)

type Request struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Response struct {
	resp.Response
	ID int64 `json:"id,omitempty"`
}

type UserRegistrar interface {
	IsUserExist(ctx context.Context, login string) (bool, error)
	SaveUser(ctx context.Context, login string, password string, date_created string) (int64, error)
}

// @Summary     Register
// @Tags        auth
// @Description register
// @ID          register
// @Accept      json
// @Produde     json
// @Param       input    body Request true "account info"
// @Router      /auth/register [post]
func New(ctx context.Context, log *slog.Logger, userRegistrar UserRegistrar) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.register.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))
		req.Login = strings.TrimSpace(req.Login)
		req.Password = strings.TrimSpace(req.Password)

		if err := validator.ValidateLogin(req.Login); err != nil {
			log.Error("invalid request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		// проверяем, существует ли пользователь с логином, переданным в запросе
		exist, err := userRegistrar.IsUserExist(ctx, req.Login)
		if err != nil {
			log.Error("failed to check if user exists", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to check if user exists"))

			return
		}

		if exist {
			log.Error("invalid request", sl.Err(fmt.Errorf("user already exists")))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("user already exists"))

			return
		}

		if err := validator.ValidatePassword(req.Password); err != nil {
			log.Error("invalid request", sl.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		hashedPassword, err := password.EncryptPassword(req.Password)
		if err != nil {
			log.Error("failed to encrypt password", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to encrypt password"))

			return
		}

		date_registered := time.Now().Format("2006-01-02 15:04:05")

		id, err := userRegistrar.SaveUser(ctx, req.Login, hashedPassword, date_registered)
		if err != nil {
			log.Error("failed to save user", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to save user"))

			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
			ID:       id,
		})
	}
}
