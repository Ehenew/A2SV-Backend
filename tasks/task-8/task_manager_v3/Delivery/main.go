package main

import (
	"context"
	"log"
	"os"
	"time"

	"a2sv-backend/task_manager_v3/Delivery/controllers"
	"a2sv-backend/task_manager_v3/Delivery/routers"
	"a2sv-backend/task_manager_v3/Infrastructure"
	"a2sv-backend/task_manager_v3/Repositories"
	"a2sv-backend/task_manager_v3/Usecases"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it — continuing using environment variables")
	}

	// MongoDB Connection
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB Atlas!")
	db := client.Database("task_manager")

	// Initialize Repositories
	userRepo := Repositories.NewMongoUserRepository(db)
	taskRepo := Repositories.NewMongoTaskRepository(db)

	// Initialize Infrastructure Services
	passwordService := Infrastructure.NewBcryptPasswordService()
	jwtService := Infrastructure.NewJWTService()

	// Initialize Usecases
	userUsecase := Usecases.NewUserUsecase(userRepo, passwordService, jwtService)
	taskUsecase := Usecases.NewTaskUsecase(taskRepo)

	// Initialize Controllers
	userController := controllers.NewUserController(userUsecase)
	taskController := controllers.NewTaskController(taskUsecase)

	// Setup Router
	r := routers.SetupRouter(taskController, userController, jwtService)

	// Run Server
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	if port[0] != ':' {
		port = ":" + port
	}

	if err := r.Run(port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
