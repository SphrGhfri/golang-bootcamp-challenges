package library

import "strings"

type Library struct {
	Capacity     int
	CurrentCount int
	Books        map[string]Book
}

type Book struct {
	Borrower string
}

func NewLibrary(capacity int) *Library {
	return &Library{
		Capacity:     capacity,
		CurrentCount: 0,
		Books:        make(map[string]Book),
	}
}

func (library *Library) AddBook(name string) string {
	lowerName := strings.ToLower(name)

	if _, exists := library.Books[lowerName]; exists {
		return "The book is already in the library"
	}

	if library.CurrentCount >= library.Capacity {
		return "Not enough capacity"
	}

	library.Books[lowerName] = Book{}
	library.CurrentCount++
	return "OK"
}

func (library *Library) BorrowBook(bookName, personName string) string {
	lowerName := strings.ToLower(bookName)

	book, exists := library.Books[lowerName]
	if !exists {
		return "The book is not defined in the library"
	}

	if book.Borrower != "" {
		return "The book is already borrowed by " + book.Borrower
	}

	book.Borrower = personName
	library.Books[lowerName] = book
	return "OK"
}

func (library *Library) ReturnBook(bookName string) string {
	lowerName := strings.ToLower(bookName)

	book, exists := library.Books[lowerName]
	if !exists {
		return "The book is not defined in the library"
	}

	if book.Borrower == "" {
		return "The book has not been borrowed"
	}

	book.Borrower = ""
	library.Books[lowerName] = book
	return "OK"
}
