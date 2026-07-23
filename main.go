package main

import (
	"database/sql"
	"log"
	"net"
	"simple_bank/api"
	db "simple_bank/db/sqlc"
	"simple_bank/gapi"
	"simple_bank/pb"
	"simple_bank/utils"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	runGrpcServer(&config, store)
	// runGinServer(&config, store)

}
func runGinServer(config *utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start(config.HttpServerAddress)

	if err != nil {
		log.Fatal("start server error")
	}
}

func runGrpcServer(config *utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterSimpleBankServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("create grpc listener error")
	}

	log.Printf("grpc server start at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("grpc server start error")
	}

}
