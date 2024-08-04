package posts

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	_ "github.com/solumD/go-blog-api/internal/http-server/models"
	resp "github.com/solumD/go-blog-api/internal/lib/api/response"
	"github.com/solumD/go-blog-api/internal/lib/logger/sl"
	"github.com/solumD/go-blog-api/internal/types"
)

type Response struct {
	resp.Response
	types.UsersPosts
	Message string `json:"message,omitempty"`
}

type PostsGetter interface {
	IsUserExist(ctx context.Context, login string) (bool, error)
	GetPosts(ctx context.Context, created_by string) (*types.UsersPosts, error)
}

// @Summary     Get posts
// @Tags        user
// @Description get posts of a user
// @ID          get
// @Accept      json
// @Produde     json
// @Param       user    path     string true "username of a user"
// @Success     200     {object} models.PostsSuccess
// @Failure     400,500 {object} models.PostsError
// @Router      /user/{user} [get]
func New(ctx context.Context, log *slog.Logger, postsGetter PostsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.post.posts.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// получаем логин пользователя из параметров запроса
		login := strings.TrimSpace(chi.URLParam(r, "login"))

		exist, err := postsGetter.IsUserExist(ctx, login)
		if err != nil {
			log.Error("failed to check if user exists", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to check if user exists"))

			return
		}
		if !exist {
			log.Error("invalid request", sl.Err(fmt.Errorf("user doesn't exist: %s", login)))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, resp.Error("user doesn't exist"))

			return
		}

		// получаем посты пользователя
		posts, err := postsGetter.GetPosts(ctx, login)
		if err != nil {
			log.Error("failed to get users's posts", sl.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to get users's posts"))

			return
		}

		// у пользователя пока нет постов
		if len(posts.Posts) == 0 {
			log.Info("user hasn't posted something yet", slog.String("user", login))

			render.JSON(w, r, Response{
				Response: resp.OK(),
				Message:  "user hasn't posted something yet",
			})

			return
		}

		log.Info("user's posts got", slog.String("user", login))

		render.JSON(w, r, Response{
			Response:   resp.OK(),
			UsersPosts: *posts,
		})
	}
}
