package like

import (
	"context"
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

type PostLiker interface {
	IsPostExist(ctx context.Context, id int) (bool, error)
	IsPostLikedByUser(ctx context.Context, id int, liked_by string) (bool, error)
	LikePost(ctx context.Context, id int, liked_by string) error
}

// @Summary Like
// @Security ApiKeyAuth
// @Tags post
// @Description like post
// @ID like
// @Accept json
// @Produde json
// @Param input body Request true "id of post to be liked"
// @Response {object} Response
// @Router /post/like [put]

func New(ctx context.Context, log *slog.Logger, postLiker PostLiker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.post.like.New"

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

		exist, err := postLiker.IsPostExist(ctx, req.ID)
		if err != nil {
			log.Error("failed to check if post exists", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to check if post exists"))

			return
		}

		if !exist {
			log.Error("invalid request", sl.Err(fmt.Errorf("post doesn't exist")))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("post doesn't exist"))

			return
		}

		login := r.Header.Get("login")
		liked, err := postLiker.IsPostLikedByUser(ctx, req.ID, login)
		if err != nil {
			log.Error("failed to check if post liked by user", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to check if post liked by user"))

			return
		}

		if liked {
			log.Error("invalid request", sl.Err(fmt.Errorf("you have already liked post %d", req.ID)))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error(fmt.Sprintf("you have already liked post %d", req.ID)))

			return
		}

		err = postLiker.LikePost(ctx, req.ID, login)
		if err != nil {
			log.Error("failed to like post", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to like post"))

			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
