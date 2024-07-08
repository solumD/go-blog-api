package posts

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
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
	IsUserExist(login string) (bool, error)
	GetPosts(created_by string) (*types.UsersPosts, error)
}

func New(log *slog.Logger, postsGetter PostsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.post.posts.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		login := strings.TrimSpace(chi.URLParam(r, "login"))

		exist, err := postsGetter.IsUserExist(login)
		if err != nil {
			log.Error("failed to check if user exists", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to check if user exists"))

			return
		}
		if !exist {
			log.Error("invalid request", sl.Err(fmt.Errorf("user doesn't exist: %s", login)))

			render.JSON(w, r, resp.Error("user doesn't exist"))

			return
		}

		posts, err := postsGetter.GetPosts(login)
		if err != nil {
			log.Error("failed to get users's posts", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to get users's posts"))

			return
		}
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