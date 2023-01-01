package api

import (
	"fmt"
	token "github.com/EliriaT/SchoolAppApi/api/token"
	"github.com/EliriaT/SchoolAppApi/config"
	"github.com/EliriaT/SchoolAppApi/service"
	"github.com/gin-gonic/gin"
)

// Serves for HTTP requests
type Server struct {
	service    service.Service
	tokenMaker token.TokenMaker
	router     *gin.Engine
	config     config.Config
}

func NewServer(service service.Service, config config.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		service:    service,
		tokenMaker: tokenMaker,
		config:     config,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", authMiddleware(server.tokenMaker), server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/schools").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("", server.createSchool)
	authRoutes.GET("/:id", server.getSchoolbyId)
	authRoutes.GET("", server.listSchools)

	server.router = router
}

// Starts the HTTP server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
