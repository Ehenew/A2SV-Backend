package services

import (
	"errors"
	"sync"
	"time"

	"library_management/models"
)

// LibraryManager defines the interface for library operations
type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	ReserveBook(bookID int, memberID int) error
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

type Library struct {
	Books   map[int]models.Book
	Members map[int]models.Member
	// Reservations maps bookID -> memberID
	Reservations map[int]int
	mu           sync.Mutex
}

func NewLibrary() *Library {
	return &Library{
		Books:        make(map[int]models.Book),
		Members:      make(map[int]models.Member),
		Reservations: make(map[int]int),
	}
}

func (l *Library) AddBook(book models.Book) {
	l.Books[book.ID] = book
}

func (l *Library) RemoveBook(bookID int) {
	delete(l.Books, bookID)
}

func (l *Library) BorrowBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, bookExists := l.Books[bookID]
	if !bookExists {
		return errors.New("book not found")
	}

	if book.Status == "Borrowed" {
		return errors.New("book is already borrowed")
	}

	// If reserved, only the reserving member can borrow
	if book.Status == "Reserved" {
		reserver, ok := l.Reservations[bookID]
		if !ok || reserver != memberID {
			return errors.New("book is reserved by another member")
		}
		// reservation consumed
		delete(l.Reservations, bookID)
	}

	member, memberExists := l.Members[memberID]
	if !memberExists {
		return errors.New("member not found")
	}

	book.Status = "Borrowed"
	l.Books[bookID] = book

	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member

	return nil
}

// ReserveBook reserves a book for a member. If reservation succeeds, an auto-cancel timer
// will unreserve the book after 5 seconds if it is not borrowed.
func (l *Library) ReserveBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, bookExists := l.Books[bookID]
	if !bookExists {
		return errors.New("book not found")
	}

	if book.Status == "Borrowed" {
		return errors.New("book is already borrowed")
	}

	if book.Status == "Reserved" {
		return errors.New("book is already reserved")
	}

	// mark reserved
	l.Reservations[bookID] = memberID
	book.Status = "Reserved"
	l.Books[bookID] = book

	// start auto-cancel timer (non-blocking)
	go func(bid int) {
		timer := time.NewTimer(5 * time.Second)
		<-timer.C
		l.mu.Lock()
		defer l.mu.Unlock()
		// if still reserved, cancel reservation
		if _, ok := l.Reservations[bid]; ok {
			delete(l.Reservations, bid)
			b := l.Books[bid]
			if b.Status == "Reserved" {
				b.Status = "Available"
				l.Books[bid] = b
			}
		}
	}(bookID)

	return nil
}

func (l *Library) ReturnBook(bookID int, memberID int) error {
	book, bookExists := l.Books[bookID]
	if !bookExists {
		return errors.New("book not found")
	}

	member, memberExists := l.Members[memberID]
	if !memberExists {
		return errors.New("member not found")
	}

	book.Status = "Available"
	l.Books[bookID] = book

	var updatedBooks []models.Book
	for _, b := range member.BorrowedBooks {
		if b.ID != bookID {
			updatedBooks = append(updatedBooks, b)
		}
	}
	member.BorrowedBooks = updatedBooks
	l.Members[memberID] = member

	return nil
}

func (l *Library) ListAvailableBooks() []models.Book {
	var availableBooks []models.Book
	for _, book := range l.Books {
		if book.Status == "Available" {
			availableBooks = append(availableBooks, book)
		}
	}
	return availableBooks
}

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	member, memberExists := l.Members[memberID]
	if !memberExists {
		return []models.Book{}
	}
	return member.BorrowedBooks
}
