package update

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/solumD/go-blog-api/internal/lib/api/response"
	"github.com/solumD/go-blog-api/internal/lib/logger/sl"
)

type Request struct {
	ID    int    `json:"id"`
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
}

type Response struct {
	resp.Response
}

type PostUpdater interface {
	GetPostCreator(ctx context.Context, id int) (string, error)
	UpdatePostTitle(ctx context.Context, id int, title string, date_updated string) error
	UpdatePostText(ctx context.Context, id int, text string, date_updated string) error
}

// @Summary Update
// @Security ApiKeyAuth
// @Tags post
// @Description update post
// @ID update
// @Accept json
// @Produde json
// @Param input body Request true "id of a post and info to be updated"
// @Response {object} Response
// @Router /post/update [patch]

func New(ctx context.Context, log *slog.Logger, PostUpdater PostUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.post.update.New"

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

		req.Title = strings.TrimSpace(req.Title)
		req.Text = strings.TrimSpace(req.Text)
		log.Info("request body decoded", slog.Any("request", req))

		if len(req.Text) == 0 && len(req.Title) == 0 {
			log.Error("invalid request", sl.Err(fmt.Errorf("title or text must be filled in")))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("title or text must be filled in"))

			return
		}

		created_by, err := PostUpdater.GetPostCreator(ctx, req.ID)
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

		date_updated := time.Now().Format("2006-01-02 15:04:05")

		if len(req.Title) > 0 {
			err := PostUpdater.UpdatePostTitle(ctx, req.ID, req.Title, date_updated)
			if err != nil {
				log.Error("failed to update post title", sl.Err(err))

				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, resp.Error("failed to update post title"))

				return
			}
		}

		if len(req.Text) > 0 {
			err := PostUpdater.UpdatePostText(ctx, req.ID, req.Text, date_updated)
			if err != nil {
				log.Error("failed to update post text", sl.Err(err))

				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, resp.Error("failed to update post text"))

				return
			}
		}

		log.Info("post updated", slog.Int("id", req.ID))

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
