package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/db/util"
	"github.com/techschool/simplebank/token"
)

//It serves HTTP requests for our banking service
type Server struct {
	config util.Config
	store db.Store
	tokenMaker token.Maker
	router *gin.Engine
}

//setup a new HHTP server & setup routing.
func NewServer(config util.Config,store db.Store) (*Server,error) {

	token,err := token.NewPasetoMaker(config.TokenSymmetricKey)  //token,err := token.NewJWTMaker(config.TokenSymmetricKey) - for jwt token
	if err !=nil{
		return nil,fmt.Errorf("cannot create token maker: %w",err)
	}
	server := &Server{
		config: config,
		store: store,
		tokenMaker: token,
	}

	if v,ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency",validCurrency) //currency is custom validator, which is used inside a json tag
	}

	server.setupRouter()
	return server,nil
}

func(server *Server) setupRouter(){
	router := gin.Default()

	router.POST("/users",server.createUser)
	router.POST("/users/login",server.loginUser)
	router.POST("/tokens/renew_access",server.RenewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts",server.createAccount)
	authRoutes.GET("/accounts/:id",server.getAccount) // : is uesd to mention that ID is a uri parameter
	authRoutes.GET("/accounts",server.listAllAccounts)
	authRoutes.DELETE("/accounts/:id", server.deleteAccount)	
	authRoutes.POST("/transfers",server.createTransfer)

	server.router = router
}
//start runs the http address on the specific address
func(server *Server) Start(address string) error{
	return server.router.Run(address)
}
func errorResponse(err error) gin.H{
	return gin.H{"error":err.Error()}
}