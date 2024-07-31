package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/solumD/go-blog-api/internal/config"
	"github.com/solumD/go-blog-api/internal/http-server/handlers/post/posts"
	"github.com/solumD/go-blog-api/internal/http-server/handlers/post/remove"
	"github.com/solumD/go-blog-api/internal/http-server/handlers/post/save"
	"github.com/solumD/go-blog-api/internal/http-server/handlers/user/login"
	"github.com/solumD/go-blog-api/internal/http-server/handlers/user/register"
	mwAuth "github.com/solumD/go-blog-api/internal/http-server/middleware/auth"
	mwLogger "github.com/solumD/go-blog-api/internal/http-server/middleware/logger"
	"github.com/solumD/go-blog-api/internal/lib/logger/sl"
	sqlite "github.com/solumD/go-blog-api/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	// инициализируем логгер
	log := InitLogger(cfg.Env)

	log.Info("starting blog-api", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	// инициализируем хранилище
	if err = storage.Init(context.TODO()); err != nil {
		log.Error("failed to init tables in storage", sl.Err(err))
		os.Exit(1)
	}

	log.Info("connected to storage")

	// инициализируем роутер
	router := chi.NewRouter()

	// инициализируем middleware логгера
	router.Use(mwLogger.New(log))

	// иниализируем вспомогательные middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// обработчики, связанные с постами
	router.Route("/post", func(r chi.Router) {
		r.Use(mwAuth.New(cfg.TokenSecret, log))
		r.Post("/create", save.New(context.Background(), log, storage))
		r.Delete("/delete", remove.New(context.Background(), log, storage))
	})

	// обработчики, связанные с пользователями
	router.Get("/user/{login}", posts.New(context.Background(), log, storage))
	router.Post("/register", register.New(context.Background(), log, storage))
	router.Post("/login", login.New(context.Background(), cfg.TokenSecret, log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func InitLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo},
			),
		)
	}
	return log
}
