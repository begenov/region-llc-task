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

func (s *Server) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		s.initLoadRoutes(v1)
	}

}
