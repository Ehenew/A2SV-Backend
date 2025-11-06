package concurrency

import (
	"library_management/services"
)

// ReservationRequest represents a request to reserve a book
type ReservationRequest struct {
	BookID   int
	MemberID int
	Resp     chan error
}

// StartReservationWorker listens on the requests channel and processes reservations concurrently.
// It spawns a goroutine per request to call ReserveBook and returns the result on the Resp channel.
func StartReservationWorker(lib services.LibraryManager, requests <-chan ReservationRequest) {
	for req := range requests {
		r := req
		go func() {
			err := lib.ReserveBook(r.BookID, r.MemberID)
			// send back the result (non-blocking if caller provided buffered channel)
			r.Resp <- err
		}()
	}
}
