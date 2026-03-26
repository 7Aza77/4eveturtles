package main

import (
	"context"
	"goevent/internal/handler"
	"goevent/internal/repository"
	"goevent/internal/usecase"
	"goevent/pkg/auth"
	"log"
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
	// 1. Подключаемся к БД
	db, err := repository.NewPostgresDB("localhost", "5432", "user", "password", "goevent")
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %s", err.Error())
	}
	log.Println("Успешное подключение к базе данных!")

	// 2. Инициализация зависимостей (DI)
	authRepo := repository.NewAuthPostgres(db)
	eventRepo := repository.NewEventPostgres(db)
	regRepo := repository.NewRegistrationPostgres(db)

	tokenManager, err := auth.NewManager("super-secret-key") // В реальности берем из ENV
	if err != nil {
		log.Fatal(err)
	}

	authUseCase := usecase.NewAuth(authRepo, tokenManager, time.Hour*12)
	eventUseCase := usecase.NewEvent(eventRepo)
	regUseCase := usecase.NewRegistration(regRepo, eventRepo)

	// 3. Redis для Rate Limiting
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // В реальности из ENV
	})

	h := handler.NewHandler(authUseCase, eventUseCase, regUseCase, tokenManager)

	// 4. Запуск сервера
	r := h.InitRouter(rdb)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 4. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
