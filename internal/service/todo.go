package service

import (
	"context"

	"github.com/begenov/region-llc-task/internal/domain"
	"github.com/begenov/region-llc-task/internal/repository"
	"github.com/begenov/region-llc-task/pkg/logger"
)

type TodoService struct {
	todoRepo repository.Todo
	userRepo repository.Users
}

func NewTodoService(todoRepo repository.Todo, userRepo repository.Users) *TodoService {
	return &TodoService{
		todoRepo: todoRepo,
		userRepo: userRepo,
	}
}

func (s *TodoService) CreateTodo(ctx context.Context, todo domain.Todo) (domain.Todo, error) {
	user, err := s.userRepo.GetUserByID(ctx, todo.UserID)
	if err != nil {
		logger.Errorf("s.userRepo.GetUserByID(): %v", err)
		return domain.Todo{}, err
	}

	todo.Actor = user.UserName

	todo, err = s.todoRepo.Create(ctx, todo)
	if err != nil {
		logger.Errorf("s.todoRepo.Create(): %v", err)
		return domain.Todo{}, err
	}

	return todo, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, todo domain.Todo) (domain.Todo, error) {
	return domain.Todo{}, nil
}
