package main

import (
	"database/sql"
	"fmt"
	"log"

	"todolist/storage"

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

	userStorage := storage.NewUserStorage(db)
	userData, err := userStorage.GetUserByID(1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(userData)

}
