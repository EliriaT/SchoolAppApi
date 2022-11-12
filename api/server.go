package api

import (
	"fmt"
	token2 "github.com/EliriaT/SchoolAppApi/api/token"
	"github.com/EliriaT/SchoolAppApi/config"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Serves for HTTP requests
type Server struct {
	store      db.Store
	tokenMaker token2.TokenMaker
	router     *gin.Engine
	config     config.Config
}

func NewServer(store db.Store, config config.Config) (*Server, error) {
	tokenMaker, err := token2.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/schools", server.createSchool)
	authRoutes.GET("/schools/:id", server.getSchoolbyId)
	authRoutes.GET("/schools", server.listSchools)

	server.router = router
}

// Starts the HTTP server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
