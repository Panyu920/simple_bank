package gapi

import (
	"fmt"
	"simple_bank/pb"
	"simple_bank/token"
	"simple_bank/utils"

	"github.com/gin-gonic/gin"
	db "simple_bank/db/sqlc"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	store      db.Store
	tokenMaker token.Maker
	config     *utils.Config
}

func NewServer(config *utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("create token maker error %v", err)
	}

	server := Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return &server, nil
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
