package main

import (
	"goevent/internal/handler"    // Твой путь из go.mod
	"goevent/internal/repository" // Твой путь из go.mod
	"log"
)

func main() {
	// 1. Подключаемся к БД (пока пропишем данные прямо тут)
	db, err := repository.NewPostgresDB("localhost", "5432", "user", "password", "goevent")
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %s", err.Error())
	}
	log.Println("Успешное подключение к базе данных!")

	// Чтобы переменная db не "ругалась", что она не используется
	_ = db

	// 2. Запуск сервера
	r := handler.InitRouter()
	err = r.Run(":8080")
	if err != nil {
		return
	}
}
