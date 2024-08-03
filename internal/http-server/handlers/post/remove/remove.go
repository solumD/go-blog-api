package remove

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

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
	GetPostCreator(ctx context.Context, id int) (string, error)
	RemovePost(ctx context.Context, id int) error
}

// @Summary     Delete
// @Security    ApiKeyAuth
// @Tags        post
// @Description delete post
// @ID          delete
// @Accept      json
// @Produde     json
// @Param       input    body Request true "id of a post to be deleted"
// @Router      /post/delete [delete]
func New(ctx context.Context, log *slog.Logger, postRemover PostRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.post.remove.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		created_by, err := postRemover.GetPostCreator(ctx, req.ID)
		if err == sql.ErrNoRows {
			log.Error("invalid request", sl.Err(fmt.Errorf("post doesn't exist: %d", req.ID)))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("post doesn't exist"))

			return
		} else if err != nil {
			log.Error("failed to check if post exists", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to check if post exists"))

			return
		}

		login := r.Header.Get("login")
		if created_by != login {
			log.Error("invalid request", sl.Err(fmt.Errorf("invalid user: %s", login)))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid user"))

			return
		}

		err = postRemover.RemovePost(ctx, req.ID)
		if err != nil {
			log.Error("failed to remove post", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to remove post"))

			return
		}

		log.Info("post removed", slog.Int("id", req.ID))

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
