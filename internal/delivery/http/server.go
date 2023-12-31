package http

import (
	_ "github.com/begenov/region-llc-task/docs"
	v1 "github.com/begenov/region-llc-task/internal/delivery/http/v1"
	"github.com/begenov/region-llc-task/internal/service"
	"github.com/begenov/region-llc-task/pkg/auth"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	engine       *gin.Engine
	userService  service.Users
	todoService  service.Todo
	tokenManager auth.TokenManager
}

func NewServer(userService service.Users, todoService service.Todo, tokenManager auth.TokenManager) *Server {
	return &Server{
		engine:       gin.New(),
		userService:  userService,
		todoService:  todoService,
		tokenManager: tokenManager,
	}
}

func (s *Server) Init(port string) error {
	s.engine.Use(logger.SetLogger())
	s.engine.Use(gin.Recovery())
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := s.engine.Group("/api")

	v1 := v1.NewServer(s.userService, s.todoService, s.tokenManager)

	v1.Init(api)

	return s.engine.Run(port)
}
