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
	router.POST("/users/twofactor", authMiddleware(server.tokenMaker), server.twoFactorLoginUser)

	schoolRoutes := router.Group("/schools").Use(authMiddleware(server.tokenMaker))

	schoolRoutes.POST("", server.createSchool)
	schoolRoutes.GET("/:id", server.getSchoolbyId)
	schoolRoutes.GET("", server.listSchools)

	classRoutes := router.Group("/class").Use(authMiddleware(server.tokenMaker))

	classRoutes.POST("", server.createClass)
	classRoutes.GET("/:id", server.getClassbyId)
	// here list of pupils can be received
	classRoutes.GET("", server.getClass)
	classRoutes.PUT("", server.changeHeadTeacherClass)

	semesterRoutes := router.Group("/semester").Use(authMiddleware(server.tokenMaker))

	semesterRoutes.POST("", server.createSemester)
	semesterRoutes.GET("", server.getSemesters)
	semesterRoutes.GET("/current", server.getCurrentSemester)

	courseRoutes := router.Group("/course").Use(authMiddleware(server.tokenMaker))

	//works
	courseRoutes.POST("", server.createCourse)
	//works
	courseRoutes.GET("", server.getCourses)
	courseRoutes.GET("/:id", server.getCourseByID)
	//works
	courseRoutes.PUT("", server.changeCourse)

	lessonRoutes := router.Group("/lesson").Use(authMiddleware(server.tokenMaker))
	//works
	lessonRoutes.POST("", server.createLesson)
	//works ideal
	lessonRoutes.GET("", server.getLessons)
	//works
	lessonRoutes.GET("/course/:id", server.getCourseLessonsByCourseID)
	//TODO CORRECT DATES UPDATE
	lessonRoutes.PUT("", server.changeLesson)

	markRoutes := router.Group("/mark").Use(authMiddleware(server.tokenMaker))
	//works
	markRoutes.POST("", server.createMark)
	//works
	markRoutes.PUT("", server.changeMark)
	//works
	markRoutes.DELETE(":id", server.deleteMark)

	//works
	router.GET("/roles", authMiddleware(server.tokenMaker), server.getRoles)
	//works
	router.GET("/teachers", authMiddleware(server.tokenMaker), server.getTeacher)

	server.router = router
}

// Starts the HTTP server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
