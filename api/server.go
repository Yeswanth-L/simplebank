package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/techschool/simplebank/db/sqlc"
)

//It serves HTTP requests for our banking service
type Server struct {
	store db.Store
	router *gin.Engine
}

//setup a new HHTP server & setup routing.
func NewServer(store db.Store) *Server {

	server := &Server{store: store}
	router := gin.Default()

	if v,ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency",validCurrency)
	}

	router.POST("/accounts",server.createAccount)
	router.GET("/accounts/:id",server.getAccount) // : is uesd to mention that ID is a uri parameter
	router.GET("/accounts",server.listAllAccounts)
	router.DELETE("/accounts/:id", server.deleteAccount)
	
	router.POST("/transfers",server.createTransfer)

	server.router = router
	return server
}

//start runs the http address on the specific address
func(server *Server) Start(address string) error{
	return server.router.Run(address)
}
func errorResponse(err error) gin.H{
	return gin.H{"error":err.Error()}
}