package main

import (
	"database/sql"
	"github.com/EliriaT/SchoolAppApi/api"
	"github.com/EliriaT/SchoolAppApi/config"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/EliriaT/SchoolAppApi/service"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	configSet, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("can not open config file")
	}

	conn, err := sql.Open(configSet.DBdriver, configSet.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)

	serverService, err := service.NewServerService(store)
	if err != nil {
		log.Fatal("cannot create create service", err)
	}

	err = serverService.CreateAdmin()
	if err != nil {
		log.Fatal("cannot create first user ", err)
	}

	server, err := api.NewServer(serverService, configSet)
	if err != nil {
		log.Fatal("cannot create new server: ", err)
	}

	//test achievement
	err = server.Start(configSet.ServerAddress)

	if err != nil {
		log.Fatal("server can not be started. ", err)
	}
}
