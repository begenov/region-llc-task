package mongo

import (
	"testing"
	"time"

	"github.com/begenov/region-llc-task/internal/domain"
	"github.com/begenov/region-llc-task/pkg/utils"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createTodo(t *testing.T) domain.Todo {
	user := createUser(t)

	todo := domain.Todo{
		UserID:   user.ID,
		Title:    utils.RandomString(15),
		ActiveAt: time.Now().Format("2006-01-02"),
		Author:   user.UserName,
		Status:   domain.Active,
	}

	newTodo, err := todoRepo.Create(ctx, todo)
	require.NoError(t, err)
	require.NotEmpty(t, newTodo)

	require.Equal(t, todo.ActiveAt, newTodo.ActiveAt)
	require.Equal(t, todo.Title, newTodo.Title)
	require.Equal(t, todo.UserID, newTodo.UserID)
	require.Equal(t, todo.Author, newTodo.Author)
	require.Equal(t, todo.Status, newTodo.Status)

	require.NotEmpty(t, newTodo.ID)
	newTodo.TodoID = newTodo.ID.Hex()

	err = todoRepo.UpdateTodoID(ctx, newTodo)
	require.NoError(t, err)

	return newTodo
}

func createTodos(t *testing.T) primitive.ObjectID {
	user := createUser(t)

	count := 10
	var todos []domain.Todo
	for i := 0; i < count; i++ {
		todo := domain.Todo{
			UserID:   user.ID,
			Title:    utils.RandomString(15),
			ActiveAt: time.Now().Format("2006-01-02"),
			Author:   user.UserName,
			Status:   domain.Active,
		}
		newTodo, err := todoRepo.Create(ctx, todo)
		require.NoError(t, err)
		require.NotEmpty(t, newTodo)

		require.Equal(t, todo.ActiveAt, newTodo.ActiveAt)
		require.Equal(t, todo.Title, newTodo.Title)
		require.Equal(t, todo.UserID, newTodo.UserID)
		require.Equal(t, todo.Author, newTodo.Author)
		require.Equal(t, todo.Status, newTodo.Status)

		require.NotEmpty(t, newTodo.ID)
		newTodo.TodoID = newTodo.ID.Hex()

		err = todoRepo.UpdateTodoID(ctx, newTodo)
		require.NoError(t, err)

		todos = append(todos, todo)
	}

	return user.ID
}

func TestTodoRepo_Create(t *testing.T) {
	createTodo(t)
}

func TestTodoRepo_GetCountByTitle(t *testing.T) {
	todo := createTodo(t)

	count, err := todoRepo.GetCountByTitle(ctx, todo.Title, todo.UserID)
	require.NoError(t, err)

	require.NotEmpty(t, count)
}

func TestTodoRepo_GetTodoByID(t *testing.T) {
	todo := createTodo(t)
	require.NotEmpty(t, todo)

	todoI, err := todoRepo.GetTodoByID(ctx, todo.ID)
	require.NoError(t, err)
	require.NotEmpty(t, todoI)

	require.Equal(t, todo.ID, todoI.ID)
	require.Equal(t, todo.UserID, todoI.UserID)
	require.Equal(t, todo.TodoID, todoI.TodoID)
	require.Equal(t, todo.ActiveAt, todoI.ActiveAt)
	require.Equal(t, todo.Title, todoI.Title)
	require.Equal(t, todo.ActiveAt, todoI.ActiveAt)
	require.Equal(t, todo.Status, todoI.Status)

	todo, err = todoRepo.GetTodoByID(ctx, primitive.NewObjectID())
	require.Error(t, err)
	require.Equal(t, err, domain.ErrNotFound)

	require.Empty(t, todo)
}

func TestTodoRepo_UpdateTodo(t *testing.T) {
	todo := createTodo(t)

	todo.Title = utils.RandomString(10)
	todo.ActiveAt = time.Now().Format("2006-01-02")

	err := todoRepo.UpdateTodo(ctx, todo)
	require.NoError(t, err)
}

func TestTodoRepo_DeleteTodoByID(t *testing.T) {
	todo := createTodo(t)

	err := todoRepo.DeleteTodoByID(ctx, todo.ID)
	require.NoError(t, err)
}

func TestTodoRepo_UpdateTodoDoneByID(t *testing.T) {
	todo := createTodo(t)

	todoU, err := todoRepo.UpdateTodoDoneByID(ctx, todo.ID, todo.UserID)
	require.NoError(t, err)
	require.NotEmpty(t, todoU)

	require.Equal(t, todo.ID, todoU.ID)
	require.Equal(t, todo.UserID, todoU.UserID)
	require.Equal(t, todo.Title, todoU.Title)
	require.Equal(t, todo.ActiveAt, todoU.ActiveAt)
	require.Equal(t, todo.Author, todoU.Author)
}

func TestTodoRepo_GetTodoByStatus(t *testing.T) {
	userID := createTodos(t)

	todos, err := todoRepo.GetTodoByStatus(ctx, domain.Active, userID)
	require.NoError(t, err)
	require.NotEmpty(t, todos)
}
