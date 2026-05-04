package main

import (
	"log"
	"net/http"
	"os"

	"final-project/pkg/api"
	"final-project/pkg/db"
)

func main() {
	// Инициализируем БД
	dbFile := "scheduler.db"
	if envDBFile := os.Getenv("TODO_DBFILE"); envDBFile != "" {
		dbFile = envDBFile
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer db.DB.Close()

	// Регистрируем API-обработчики
	api.Init()

	// Файловый сервер
	http.Handle("/", http.FileServer(http.Dir("./web")))

	// Определяем порт
	port := "7540"
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	log.Printf("Сервер запущен на http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
