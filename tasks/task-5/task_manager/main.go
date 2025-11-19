package main

import "a2sv-backend/task_manager/router"

func main() {
	r := router.SetupRouter()
	r.Run()
}
