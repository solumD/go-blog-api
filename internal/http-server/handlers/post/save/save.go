package save

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	_ "github.com/solumD/go-blog-api/internal/http-server/models"
	resp "github.com/solumD/go-blog-api/internal/lib/api/response"
	"github.com/solumD/go-blog-api/internal/lib/logger/sl"
)

type Request struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Response struct {
	resp.Response
	ID int64 `json:"id,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.30.0 --name=PostSaver
type PostSaver interface {
	SavePost(ctx context.Context, created_by string, title string, text string, date_created string) (int64, error)
}

// @Summary     Create
// @Security    ApiKeyAuth
// @Tags        post
// @Description create post
// @ID          create
// @Accept      json
// @Produde     json
// @Param       input   body     Request true "post info"
// @Success     200     {object} models.SaveSuccess
// @Failure     400,500 {object} models.SaveError
// @Router      /post/create [post]
func New(ctx context.Context, log *slog.Logger, postSaver PostSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.post.save.New"

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

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		req.Title = strings.TrimSpace(req.Title)
		req.Text = strings.TrimSpace(req.Text)
		log.Info("request body decoded", slog.Any("request", req))

		if len(req.Title) == 0 || len(req.Text) == 0 {
			log.Error("invalid request", sl.Err(errors.New("post's title and text can't be empty")))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("post's title and text can't be empty"))

			return
		}

		login := r.Header.Get("login")

		date_created := time.Now().Format("2006-01-02 15:04:05")

		id, err := postSaver.SavePost(ctx, login, req.Title, req.Text, date_created)
		if err != nil {
			log.Error("failed to save post", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to save post"))

			return
		}

		log.Info("post created", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			ID:       id,
		})
	}
}
