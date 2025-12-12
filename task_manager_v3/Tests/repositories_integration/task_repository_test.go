package repositories_integration_test

import (
	"a2sv-backend/task_manager_v3/Domain"
	"a2sv-backend/task_manager_v3/Repositories"
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepositorySuite struct {
	suite.Suite
	client *mongo.Client
	db     *mongo.Database
	repo   Repositories.TaskRepository
}

func (s *TaskRepositorySuite) SetupSuite() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		s.T().Skip("Could not connect to MongoDB: ", err)
		return
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		s.T().Skip("Could not ping MongoDB: ", err)
		return
	}

	s.client = client
	s.db = client.Database("test_task_manager_v3")
	s.repo = Repositories.NewMongoTaskRepository(s.db)
}

func (s *TaskRepositorySuite) TearDownSuite() {
	if s.client != nil {
		s.db.Drop(context.Background())
		s.client.Disconnect(context.Background())
	}
}

func (s *TaskRepositorySuite) SetupTest() {
	if s.db != nil {
		s.db.Collection("tasks").DeleteMany(context.Background(), bson.M{})
	}
}

func (s *TaskRepositorySuite) TestCreateTask() {
	if s.db == nil {
		s.T().Skip("MongoDB not available")
	}

	task := Domain.Task{
		Title:       "Integration Test Task",
		Description: "Description",
		DueDate:     time.Now(),
		Status:      "pending",
	}

	createdTask, err := s.repo.Create(task)

	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), createdTask.ID)
	assert.Equal(s.T(), task.Title, createdTask.Title)
}

func (s *TaskRepositorySuite) TestFindAll() {
	if s.db == nil {
		s.T().Skip("MongoDB not available")
	}

	task1 := Domain.Task{Title: "Task 1"}
	task2 := Domain.Task{Title: "Task 2"}

	s.repo.Create(task1)
	s.repo.Create(task2)

	tasks, err := s.repo.FindAll()

	assert.NoError(s.T(), err)
	assert.Len(s.T(), tasks, 2)
}

func TestTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(TaskRepositorySuite))
}
