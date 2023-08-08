package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/begenov/region-llc-task/internal/config"
	"github.com/begenov/region-llc-task/pkg/database"
	"github.com/begenov/region-llc-task/pkg/logger"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
)

var userRepo *UserRepo
var todoRepo *TodoRepo
var ctx context.Context

func init() {
	client := createTestDatabaseClient()
	db := client.Database("test")
	userRepo = NewUserRepo(db)
	todoRepo = NewTodoRepo(db)
}

func createTestDatabaseClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        "mongo:6",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Waiting for connections"),
			wait.ForListeningPort("27017/tcp"),
		),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		logger.Fatalf("testcontainers.GenericContainer(): %v", err)
	}

	ip, err := container.Host(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	port, err := container.MappedPort(ctx, "27017")
	if err != nil {
		logger.Fatal(err)
	}

	mongoURI := fmt.Sprintf("mongodb://%s:%s", ip, port.Port())

	client, err := database.NewClient(ctx, config.ConfigMongo{
		Uri: mongoURI,
	})
	if err != nil {
		logger.Fatalf("database.NewClient(): %v", err)
	}

	return client
}
