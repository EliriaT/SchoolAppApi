package api

import (
	"fmt"
	token "github.com/EliriaT/SchoolAppApi/api/token"
	"github.com/EliriaT/SchoolAppApi/config"
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Serves for HTTP requests
type Server struct {
	store      db.Store
	tokenMaker token.TokenMaker
	router     *gin.Engine
	config     config.Config
}

func NewServer(store db.Store, config config.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
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
