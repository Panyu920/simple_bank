package main

import (
	"database/sql"
	"log"
	"simple_bank/api"
	db "simple_bank/db/sqlc"
	"simple_bank/utils"

	_ "github.com/lib/pq"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("load config err ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("open db failed", err)
	}

	store := db.NewStore(conn)

	server, err := api.NewServer(&config, store)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("start server error")
	}

}
