package main

import (
	"context"
	"goevent/internal/config"
	"goevent/internal/handler"
	"goevent/internal/repository"
	"goevent/internal/usecase"
	"goevent/pkg/auth"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

// @title           GoEvent API
// @version         1.0
// @description     REST API for the GoEvent student event management platform.

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	cfg := config.MustLoad()

	var log *slog.Logger
	if cfg.Env == "local" {
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	slog.SetDefault(log)

	log.Info("starting goevent", slog.String("env", cfg.Env))

	db, err := repository.NewPostgresDB(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
	)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	log.Info("database connected")

	authRepo := repository.NewAuthPostgres(db)
	eventRepo := repository.NewEventPostgres(db)
	regRepo := repository.NewRegistrationPostgres(db)

	tokenManager, err := auth.NewManager(cfg.Auth.JWTSecret)
	if err != nil {
		log.Error("failed to init token manager", "error", err)
		os.Exit(1)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Host + ":" + cfg.Redis.Port,
	})

	authUseCase := usecase.NewAuth(authRepo, tokenManager, time.Hour*12)
	eventUseCase := usecase.NewEvent(eventRepo, rdb)
	regUseCase := usecase.NewRegistration(regRepo, eventRepo)

	h := handler.NewHandler(authUseCase, eventUseCase, regUseCase, tokenManager)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      h.InitRouter(rdb),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server listen error", "error", err)
			os.Exit(1)
		}
	}()

	log.Info("server started", slog.String("address", cfg.HTTPServer.Address))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	log.Info("server stopped gracefully")
}
