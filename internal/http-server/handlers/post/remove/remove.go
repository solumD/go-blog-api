package remove

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/solumD/go-blog-api/internal/lib/api/response"
	"github.com/solumD/go-blog-api/internal/lib/logger/sl"
)

type Request struct {
	ID int `json:"id"`
}

type Response struct {
	resp.Response
}

type PostRemover interface {
	GetPostCreator(id int) (string, error)
	RemovePost(id int) error
}

// TODO: получения логина пользователя из JWT-токена
func New(log *slog.Logger, postRemover PostRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.post.save.Delete"

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

		created_by, err := postRemover.GetPostCreator(req.ID)
		if err == sql.ErrNoRows {
			log.Error("invalid request", sl.Err(fmt.Errorf("post doesn't exist: %d", req.ID)))

			render.JSON(w, r, resp.Error("post doesn't exist"))

			return
		} else if err != nil {
			log.Error("failed to check if post exists", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to check if post exists"))

			return
		}

		login := "test"
		if created_by != login {
			log.Error("invalid request", sl.Err(fmt.Errorf("invalid user: %s", login)))

			render.JSON(w, r, resp.Error("invalid user"))

			return
		}

		err = postRemover.RemovePost(req.ID)
		if err != nil {
			log.Error("failed to remove post", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to remove post"))

			return
		}

		log.Info("post removed", slog.Int("id", req.ID))

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
