package service

import (
	"context"

	"github.com/begenov/region-llc-task/internal/domain"
	"github.com/begenov/region-llc-task/internal/repository"
	"github.com/begenov/region-llc-task/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	if err := s.validateTitle(ctx, todo.UserID, todo.Title, todo.ActiveAt); err != nil {
		return domain.Todo{}, err
	}

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

	todo.TodoID = todo.ID.Hex()

	return todo, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, todo domain.Todo) (domain.Todo, error) {

	todoID, err := primitive.ObjectIDFromHex(todo.TodoID)
	if err != nil {
		return domain.Todo{}, err
	}

	todo.ID = todoID

	if err := s.todoRepo.UpdateTodo(ctx, todo); err != nil {
		logger.Errorf("s.todoRepo.UpdateTodo(): %v", err)
		return domain.Todo{}, err
	}

	todo, err = s.todoRepo.GetTodoByID(ctx, todo.ID)
	if err != nil {
		logger.Errorf("s.todoRepo.GetTodoByID(): %v", err)
		return domain.Todo{}, err
	}

	todo.TodoID = todo.ID.Hex()

	return todo, nil
}

func (s *TodoService) DeleteTodoByID(ctx context.Context, id string) error {
	todoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if err := s.todoRepo.DeleteTodoByID(ctx, todoID); err != nil {
		return err
	}

	return nil
}

func (s *TodoService) TodoDoneUpdateByID(ctx context.Context, id string) {

}

func (s *TodoService) validateTitle(ctx context.Context, user_id primitive.ObjectID, title, activeAt string) error {

	count, err := s.todoRepo.GetCountByTitle(ctx, title, user_id)
	if err != nil {
		return err
	}

	if count > 0 {
		return domain.ErrTitleAlreadyExists
	}

	return nil
}
