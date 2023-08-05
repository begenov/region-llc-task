package service

import "github.com/begenov/region-llc-task/internal/repository"

type TodoService struct {
	todoRepo repository.Todo
}

func NewTodoService(todoRepo repository.Todo) *TodoService {
	return &TodoService{
		todoRepo: todoRepo,
	}
}
