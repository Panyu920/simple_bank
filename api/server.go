package api

import db "simple_bank/db/sqlc"
import "github.com/gin-gonic/gin"

type Server struct {
	store   *db.Store
	rounter *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := Server{store: store}

	server.rounter = gin.Default()

	server.rounter.POST("/account", server.createAccount)
	server.rounter.GET("/account/:id", server.getAccount)
	server.rounter.GET("/account/", server.listAccount)

	return &server
}

func (server *Server) Start(address string) error {
	return server.rounter.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
