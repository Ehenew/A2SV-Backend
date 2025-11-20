package main

import (
	"a2sv-backend/task_manager/data"
	"a2sv-backend/task_manager/router"
)

func main() {
	data.InitMongoDB()
	r := router.SetupRouter()
	r.Run()
}
