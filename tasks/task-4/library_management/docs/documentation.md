# Library Management System

This is a simple console-based library management system written in Go.

## How to Run

1. Navigate to the `library_management` directory.
2. Run the application using `go run main.go`.

## Concurrency (Reservations)

This project supports concurrent reservation processing using Goroutines, Channels and a Mutex to
prevent race conditions when updating shared state.

- `services.Library` holds a `Reservations` map and a `mu sync.Mutex` to protect access to `Books`,
	`Members` and `Reservations`.
- `ReserveBook(bookID, memberID)` marks a book as `Reserved` (if available) and starts an auto-cancel
	timer: if the reservation is not converted to a borrow within 5 seconds, the reservation is cancelled
	and the book becomes `Available` again.
- `concurrency.StartReservationWorker` listens on a `chan ReservationRequest` and processes each request
	concurrently in its own goroutine by calling `LibraryManager.ReserveBook`.

Use the console menu option "Reserve Book (async)" to send reservation requests to the worker via the
requests channel. The worker will process multiple reservation requests in parallel while `Library`'s
mutex ensures updates are safe.
