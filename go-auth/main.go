package main

import (
	"go-auth/routes"
)

func main() {
	router := routes.SetupRouter()
	router.Run()
}
