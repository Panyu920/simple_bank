package api

import (
	"fmt"
	"log"
	db "simple_bank/db/sqlc"
	"simple_bank/token"
	"simple_bank/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store      db.Store
	rounter    *gin.Engine
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个名为 "currency" 的自定义验证器
		v.RegisterValidation("currency", currencyValidator)
	} else {
		log.Fatal("注册 currency 验证器失败")
	}

	server.setupRouter()

	return &server, nil
}

func (server *Server) setupRouter() {
	server.rounter = gin.Default()

	server.rounter.POST("/user", server.createUser)
	server.rounter.POST("/user/login", server.loginUser)

	authRouter := server.rounter.Group("/").Use(authMiddleWare(server.tokenMaker))
	authRouter.POST("/account", server.createAccount)
	authRouter.GET("/account/:id", server.getAccount)
	authRouter.GET("/account/", server.listAccount)

	authRouter.POST("/transfer", server.createTransfer)

}

func (server *Server) Start(address string) error {
	return server.rounter.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
