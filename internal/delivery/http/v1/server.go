package v1

import (
	"github.com/begenov/region-llc-task/internal/service"
	"github.com/begenov/region-llc-task/pkg/auth"
	"github.com/gin-gonic/gin"
)

type Server struct {
	userService  service.Users
	todoService  service.Todo
	tokenManager auth.TokenManager
}

func NewServer(userService service.Users, todoService service.Todo, tokenManager auth.TokenManager) *Server {
	return &Server{
		userService:  userService,
		todoService:  todoService,
		tokenManager: tokenManager,
	}
}

func (s *Server) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		s.initLoadRoutes(v1)
	}

}
