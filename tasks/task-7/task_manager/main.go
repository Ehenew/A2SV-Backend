package main

import (
	"log"
	"os"

	"a2sv-backend/task_manager/data"
	"a2sv-backend/task_manager/router"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env if present (non-fatal)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it — continuing using environment variables")
	}

	data.InitMongoDB()
	r := router.SetupRouter()

	// Use PORT env var if set, otherwise default to :8080
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