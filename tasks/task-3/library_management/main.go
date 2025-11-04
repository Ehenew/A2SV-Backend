package main

import (
	"library_management/controllers"
	"library_management/models"
	"library_management/services"
)

func main() {
	// initialize library service
	lib := services.NewLibrary()

	// seed some data
	lib.AddBook(models.Book{ID: 1, Title: "The Go Programming Language", Author: "Donovan & Kernighan", Status: "Available"})
	lib.AddBook(models.Book{ID: 2, Title: "Clean Code", Author: "Robert C. Martin", Status: "Available"})

	// add a sample member
	lib.Members[1] = models.Member{ID: 1, Name: "Default Member"}

	// start console UI
	controllers.Start(lib)
}
