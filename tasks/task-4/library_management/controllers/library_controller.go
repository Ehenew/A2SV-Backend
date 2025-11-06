package controllers

import (
	"bufio"
	"fmt"
	"library_management/concurrency"
	"library_management/models"
	"library_management/services"
	"os"
	"strconv"
)

func Start(library services.LibraryManager, requests chan concurrency.ReservationRequest) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nLibrary Management System")
		fmt.Println("1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Borrow Book")
		fmt.Println("4. Return Book")
		fmt.Println("5. List Available Books")
		fmt.Println("6. List Borrowed Books")
		fmt.Println("7. Reserve Book (sync)")
		fmt.Println("8. Reserve Book (async)")
		fmt.Println("9. Simulate concurrent reservations")
		fmt.Println("10. Exit")
		fmt.Print("Enter choice: ")

		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			addBook(scanner, library)
		case "2":
			removeBook(scanner, library)
		case "3":
			borrowBook(scanner, library)
		case "4":
			returnBook(scanner, library)
		case "5":
			listAvailableBooks(library)
		case "6":
			listBorrowedBooks(scanner, library)
		case "7":
			reserveBook(scanner, library)
		case "8":
			asyncReserve(scanner, library, requests)
		case "9":
			simulateConcurrentReservations(library, requests)
		case "10":
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func reserveBook(scanner *bufio.Scanner, library services.LibraryManager) {
	fmt.Print("Enter book ID to reserve: ")
	scanner.Scan()
	bookID, _ := strconv.Atoi(scanner.Text())

	fmt.Print("Enter member ID: ")
	scanner.Scan()
	memberID, _ := strconv.Atoi(scanner.Text())

	err := library.ReserveBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book reserved successfully (sync).")
	}
}

func asyncReserve(scanner *bufio.Scanner, library services.LibraryManager, requests chan concurrency.ReservationRequest) {
	fmt.Print("Enter book ID to reserve: ")
	scanner.Scan()
	bookID, _ := strconv.Atoi(scanner.Text())

	fmt.Print("Enter member ID: ")
	scanner.Scan()
	memberID, _ := strconv.Atoi(scanner.Text())

	resp := make(chan error, 1)
	requests <- concurrency.ReservationRequest{BookID: bookID, MemberID: memberID, Resp: resp}
	err := <-resp
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book reserved successfully (async).")
	}
}

// simulateConcurrentReservations demonstrates handling multiple concurrent reservation requests
func simulateConcurrentReservations(library services.LibraryManager, requests chan concurrency.ReservationRequest) {
	// For the demo, attempt to reserve book ID 1 by member 1 and member 2 at the same time
	bookID := 1
	memberA := 1
	memberB := 2

	respA := make(chan error, 1)
	respB := make(chan error, 1)

	// send both requests concurrently
	go func() { requests <- concurrency.ReservationRequest{BookID: bookID, MemberID: memberA, Resp: respA} }()
	go func() { requests <- concurrency.ReservationRequest{BookID: bookID, MemberID: memberB, Resp: respB} }()

	errA := <-respA
	errB := <-respB

	fmt.Printf("Member %d reservation result: %v\n", memberA, errA)
	fmt.Printf("Member %d reservation result: %v\n", memberB, errB)
}

func addBook(scanner *bufio.Scanner, library services.LibraryManager) {
	fmt.Print("Enter book ID: ")
	scanner.Scan()
	id, _ := strconv.Atoi(scanner.Text())

	fmt.Print("Enter book title: ")
	scanner.Scan()
	title := scanner.Text()

	fmt.Print("Enter book author: ")
	scanner.Scan()
	author := scanner.Text()

	library.AddBook(models.Book{ID: id, Title: title, Author: author, Status: "Available"})
	fmt.Println("Book added successfully.")
}

func removeBook(scanner *bufio.Scanner, library services.LibraryManager) {
	fmt.Print("Enter book ID to remove: ")
	scanner.Scan()
	id, _ := strconv.Atoi(scanner.Text())
	library.RemoveBook(id)
	fmt.Println("Book removed successfully.")
}

func borrowBook(scanner *bufio.Scanner, library services.LibraryManager) {
	fmt.Print("Enter book ID to borrow: ")
	scanner.Scan()
	bookID, _ := strconv.Atoi(scanner.Text())

	fmt.Print("Enter member ID: ")
	scanner.Scan()
	memberID, _ := strconv.Atoi(scanner.Text())

	err := library.BorrowBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book borrowed successfully.")
	}
}

func returnBook(scanner *bufio.Scanner, library services.LibraryManager) {
	fmt.Print("Enter book ID to return: ")
	scanner.Scan()
	bookID, _ := strconv.Atoi(scanner.Text())

	fmt.Print("Enter member ID: ")
	scanner.Scan()
	memberID, _ := strconv.Atoi(scanner.Text())

	err := library.ReturnBook(bookID, memberID)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Book returned successfully.")
	}
}

func listAvailableBooks(library services.LibraryManager) {
	books := library.ListAvailableBooks()
	fmt.Println("\nAvailable Books:")
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
	}
}

func listBorrowedBooks(scanner *bufio.Scanner, library services.LibraryManager) {
	fmt.Print("Enter member ID: ")
	scanner.Scan()
	memberID, _ := strconv.Atoi(scanner.Text())

	books := library.ListBorrowedBooks(memberID)
	fmt.Println("\nBorrowed Books:")
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.ID, book.Title, book.Author)
	}
}
