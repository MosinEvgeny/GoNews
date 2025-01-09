package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/memdb"
	"GoNews/pkg/storage/mongo"
	"GoNews/pkg/storage/postgres"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {

	// Загружаем переменные окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Создаём объект сервера.
	var srv server

	// Выбираем тип базы данных из переменной окружения
	dbType := os.Getenv("DATABASE_TYPE")

	// Создаём объект базы данных в зависимости от выбранного типа
	var db storage.Interface
	switch dbType {
	case "postgres":
		connStr := os.Getenv("POSTGRES_CONNECTION_STRING")
		db, err = postgres.New(connStr)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to create postgres storage: %w", err))
		}
		log.Println("Using PostgreSQL database")
	case "mongo":
		connStr := os.Getenv("MONGO_CONNECTION_STRING")

		db, err = mongo.New(connStr)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to create mongo storage: %w", err))
		}
		log.Println("Using MongoDB database")
	default:
		db = memdb.New()
		log.Println("Using in-memory database")
	}

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер.
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", srv.api.Router())
}
