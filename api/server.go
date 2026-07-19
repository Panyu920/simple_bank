package api

import (
	"log"
	db "simple_bank/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store   db.Store
	rounter *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := Server{store: store}

	server.rounter = gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个名为 "currency" 的自定义验证器
		v.RegisterValidation("currency", currencyValidator)
	} else {
		log.Fatal("注册 currency 验证器失败")
	}

	server.rounter.POST("/user", server.createUser)

	server.rounter.POST("/account", server.createAccount)
	server.rounter.GET("/account/:id", server.getAccount)
	server.rounter.GET("/account/", server.listAccount)

	server.rounter.POST("/transfer", server.createTransfer)

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
