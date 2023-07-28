package main

import (
	"database/sql"
	"fmt"
	"log"
	api "todolist/internal/api"
	"todolist/internal/config"
	"todolist/internal/service"
	storage "todolist/internal/storage"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 15432
	user     = "postgres"
	password = "dev"
	dbname   = "postgres"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	redisClient, err := storage.NewRedisClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	userStorage := storage.NewUserStorage(db, redisClient)
	userService := service.NewUser(userStorage)
	userHandler := api.NewHandler(userService)
	authMW := api.NewAuthMiddleware(userService)
	server := api.NewServer("8080", userHandler, authMW)
	err = server.Run()
	if err != nil {
		log.Fatal(" server run: ", err)
	}
}
