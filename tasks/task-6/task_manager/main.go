package main

import (
	"log"

	"a2sv-backend/task_manager/data"
	"a2sv-backend/task_manager/router"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	data.InitMongoDB()
	r := router.SetupRouter()
	r.Run()
}
