package api

import (
	db "github.com/EliriaT/SchoolAppApi/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Serves for HTTP requests
type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	router.POST("/schools", server.createSchool)
	router.GET("/schools/:id", server.getSchoolbyId)
	router.GET("/schools", server.listSchools)

	server.router = router
	return server
}

// Starts the HTTP server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
