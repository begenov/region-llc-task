package http

import (
	v1 "github.com/begenov/region-llc-task/internal/delivery/http/v1"
	"github.com/begenov/region-llc-task/internal/service"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine      *gin.Engine
	userService service.Users
	todoService service.Todo
}

func NewServer(userService service.Users, todoService service.Todo) *Server {
	return &Server{
		engine:      gin.New(),
		userService: userService,
		todoService: todoService,
	}
}

func (s *Server) Init(port string) error {
	s.engine.Use(logger.SetLogger())
	s.engine.Use(gin.Recovery())

	api := s.engine.Group("/api")

	v1 := v1.NewServer(s.userService, s.todoService)

	v1.Init(api)

	return s.engine.Run(port)
}
