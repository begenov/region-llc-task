package service

import (
	"context"
	"sort"
	"time"

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

	if err := validateTodo(todo.Title, todo.ActiveAt); err != nil {
		logger.Errorf("validateTodo(): %v", err)
		return domain.Todo{}, err
	}

	if err := s.checkTitle(ctx, todo.UserID, todo.Title); err != nil {
		return domain.Todo{}, err
	}

	user, err := s.userRepo.GetUserByID(ctx, todo.UserID)
	if err != nil {
		logger.Errorf("s.userRepo.GetUserByID(): %v", err)
		return domain.Todo{}, err
	}

	todo.Author = user.UserName

	todo, err = s.todoRepo.Create(ctx, todo)
	if err != nil {
		logger.Errorf("s.todoRepo.Create(): %v", err)
		return domain.Todo{}, err
	}

	return todo, nil
}

func (s *TodoService) UpdateTodo(ctx context.Context, todo domain.Todo) (domain.Todo, error) {

	if err := validateTodo(todo.Title, todo.ActiveAt); err != nil {
		logger.Errorf("validateTodo(): %v", err)
		return domain.Todo{}, err
	}

	if err := s.checkTitle(ctx, todo.UserID, todo.Title); err != nil {
		return domain.Todo{}, err
	}

	if err := s.todoRepo.UpdateTodo(ctx, todo); err != nil {
		logger.Errorf("s.todoRepo.UpdateTodo(): %v", err)
		return domain.Todo{}, err
	}

	todo, err := s.todoRepo.GetTodoByID(ctx, todo.ID)
	if err != nil {
		logger.Errorf("s.todoRepo.GetTodoByID(): %v", err)
		return domain.Todo{}, err
	}

	return todo, nil
}

func (s *TodoService) DeleteTodoByID(ctx context.Context, id string) error {
	todoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrTodoInvalidId
	}

	if err := s.todoRepo.DeleteTodoByID(ctx, todoID); err != nil {
		return err
	}

	return nil
}

func (s *TodoService) UpdateTodoDoneByID(ctx context.Context, id string, userID primitive.ObjectID) (domain.Todo, error) {
	todoID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Todo{}, domain.ErrTodoInvalidId
	}

	todo, err := s.todoRepo.UpdateTodoDoneByID(ctx, todoID, userID)
	if err != nil {
		return domain.Todo{}, err
	}

	if todo.Status == domain.Active {
		logger.Info("TODO")
		todo.Status = domain.Done
	}

	return todo, nil
}

func (s *TodoService) GetTodosByStatus(ctx context.Context, status string, userID primitive.ObjectID) ([]domain.Todo, error) {
	if status == "" {
		status = "active"
	}

	todos, err := s.todoRepo.GetTodoByStatus(ctx, status, userID)
	if err != nil {
		return nil, err
	}

	var result []domain.Todo

	for _, todo := range todos {
		activeAtTime, err := parseTimeString(todo.ActiveAt)
		if err != nil {
			logger.Errorf("parseTimeString(): %v", err)
			continue
		}

		if status == domain.Active && activeAtTime.After(time.Now()) {
			continue
		}

		weekendTitle := ""
		if activeAtTime.Weekday() == time.Saturday || activeAtTime.Weekday() == time.Sunday {
			weekendTitle = "ВЫХОДНОЙ - "
		}

		todo.Title = weekendTitle + todo.Title
		result = append(result, todo)
	}

	sortByCreatedAt(result)

	return result, nil
}

func (s *TodoService) checkTitle(ctx context.Context, user_id primitive.ObjectID, title string) error {

	count, err := s.todoRepo.GetCountByTitle(ctx, title, user_id)
	if err != nil {
		return err
	}

	if count > 0 {
		return domain.ErrTitleAlreadyExists
	}

	return nil
}

func sortByCreatedAt(todos []domain.Todo) {
	sort.Slice(todos, func(i, j int) bool {
		timeI, err := parseTimeString(todos[i].ActiveAt)
		if err != nil {
			return true
		}

		timeJ, err := parseTimeString(todos[j].ActiveAt)
		if err != nil {
			return false
		}

		return timeI.Before(timeJ)
	})

}

func parseTimeString(activeAt string) (time.Time, error) {
	layout := "2006-01-02"

	t, err := time.Parse(layout, activeAt)
	if err != nil {
		return time.Time{}, domain.ErrIncorrectDateFormat
	}

	return t, nil
}

func validateTodo(title string, activeAt string) error {

	if title == "" {
		return domain.ErrInvalidTitle
	}

	if len(title) > 200 {
		return domain.ErrHeaderLength
	}

	t, err := time.Parse("2006-01-02", activeAt)
	if err != nil {
		return domain.ErrIncorrectDateFormat
	}

	currentTime := time.Now()

	if t.Before(currentTime) {
		return domain.ErrTodoActiveAtData
	}

	return nil
}
