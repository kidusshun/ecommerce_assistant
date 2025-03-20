package main

import (
	"database/sql"
	"log"

	"github.com/kidusshun/ecom_bot/cmd/api"
	"github.com/kidusshun/ecom_bot/config"
	"github.com/kidusshun/ecom_bot/db"
)

func main() {
	db, err := db.NewMySQLStorage(
		config.Envs.DBUser,
		config.Envs.DBPassword,
		config.Envs.DBAddress,
		config.Envs.DBName,
	)
	if err != nil {
		log.Fatal(err)
	}
	initStorage(db)
	server := api.NewAPIServer(":8080", db)
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB successfully connected")
}
