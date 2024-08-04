package unlike

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	_ "github.com/solumD/go-blog-api/internal/http-server/models"
	resp "github.com/solumD/go-blog-api/internal/lib/api/response"
	"github.com/solumD/go-blog-api/internal/lib/logger/sl"
)

type Request struct {
	ID int `json:"id"`
}

type Response struct {
	resp.Response
}

type PostUnLiker interface {
	IsPostExist(ctx context.Context, id int) (bool, error)
	IsPostLikedByUser(ctx context.Context, id int, liked_by string) (bool, error)
	UnlikePost(ctx context.Context, id int, liked_by string) error
}

// @Summary     Unlike
// @Security    ApiKeyAuth
// @Tags        post
// @Description unlike post
// @ID          unlike
// @Accept      json
// @Produde     json
// @Param       input   body     Request true "id of post to be unliked"
// @Success     200     {object} models.UnlikeSuccess
// @Failure     400,500 {object} models.UnlikeError
// @Router      /post/unlike [put]
func New(ctx context.Context, log *slog.Logger, postUnliker PostUnLiker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.post.unlike.New"

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

		exist, err := postUnliker.IsPostExist(ctx, req.ID)
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
		liked, err := postUnliker.IsPostLikedByUser(ctx, req.ID, login)
		if err != nil {
			log.Error("failed to check if post liked by user", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to check if post liked by user"))

			return
		}

		if !liked {
			log.Error("invalid request", sl.Err(fmt.Errorf("you haven't liked post %d", req.ID)))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error(fmt.Sprintf("you haven't liked post %d", req.ID)))

			return
		}

		err = postUnliker.UnlikePost(ctx, req.ID, login)
		if err != nil {
			log.Error("failed to unlike post", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to unlike post"))

			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}
