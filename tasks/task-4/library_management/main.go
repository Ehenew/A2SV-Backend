package main

import (
	"library_management/concurrency"
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

	// add sample members
	lib.Members[1] = models.Member{ID: 1, Name: "Alice"}
	lib.Members[2] = models.Member{ID: 2, Name: "Bob"}

	// create reservation request channel and start worker
	requests := make(chan concurrency.ReservationRequest)
	go concurrency.StartReservationWorker(lib, requests)

	// start console UI (pass the requests channel)
	controllers.Start(lib, requests)
}
