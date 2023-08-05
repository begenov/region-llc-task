package v1

import (
	"github.com/begenov/region-llc-task/internal/service"
	"github.com/gin-gonic/gin"
)

type Server struct {
	userService service.Users
	todoService service.Todo
}

func NewServer(userService service.Users, todoService service.Todo) *Server {
	return &Server{
		userService: userService,
		todoService: todoService,
	}
}

func (s *Server) InitRouters(api *gin.RouterGroup) {
}
