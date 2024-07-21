package main

import (
	"log"
	"os"

	"github.com/YonkLongSchlong/Todo-BE/packages/server"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}

	dbConnectString := os.Getenv("DB_CONNECT_STRING")
	db, err := sqlx.Connect("mysql", dbConnectString)
	if err != nil {
		log.Fatal(err)
	}
	server := server.NewApiServer(":8080", db)
	server.Run()
}
