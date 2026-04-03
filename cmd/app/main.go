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

// @title GoEvent API
// @version 1.0
// @description API Server for GoEvent Application

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// 1. Загрузка конфигурации
	cfg := config.MustLoad()

	// 2. Инициализация логгера (slog)
	var log *slog.Logger
	if cfg.Env == "local" {
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	slog.SetDefault(log)

	log.Info("starting goevent application", slog.String("env", cfg.Env))

	// 3. Подключение к БД
	db, err := repository.NewPostgresDB(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
	)
	if err != nil {
		log.Error("failed to connect to db", "error", err)
		os.Exit(1)
	}
	log.Debug("database connected successfully")

	// 4. Инициализация зависимостей (DI)
	authRepo := repository.NewAuthPostgres(db)
	eventRepo := repository.NewEventPostgres(db)
	regRepo := repository.NewRegistrationPostgres(db)

	tokenManager, err := auth.NewManager(cfg.Auth.JWTSecret)
	if err != nil {
		log.Error("failed to init token manager", "error", err)
		os.Exit(1)
	}

	// 5. Redis для Rate Limiting и Кэширования
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Host + ":" + cfg.Redis.Port,
	})

	authUseCase := usecase.NewAuth(authRepo, tokenManager, time.Hour*12)
	eventUseCase := usecase.NewEvent(eventRepo, rdb)
	regUseCase := usecase.NewRegistration(regRepo, eventRepo)

	h := handler.NewHandler(authUseCase, eventUseCase, regUseCase, tokenManager)

	// 6. Запуск сервера
	r := h.InitRouter(rdb)
	srv := &http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen and serve error", "error", err)
			os.Exit(1)
		}
	}()

	log.Info("server started", slog.String("address", cfg.HTTPServer.Address))

	// 7. Graceful Shutdown
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

	log.Info("server exiting")
}
