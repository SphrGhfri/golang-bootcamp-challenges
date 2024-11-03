package main

import (
	"errors"
	"fmt"
)

// Node represents a node in the doubly linked list
type Node struct {
	url  string
	prev *Node
	next *Node
}

// BrowserHistory manages the browsing history using a doubly linked list
type BrowserHistory struct {
	current *Node
	head    *Node
	tail    *Node
}

// NewBrowserHistory initializes and returns an instance of BrowserHistory
func NewBrowserHistory() *BrowserHistory {
	return &BrowserHistory{}
}

// VisitNewPage adds a new page to the history
func (bh *BrowserHistory) VisitNewPage(url string) {
	newNode := &Node{url: url}

	// If we are not at the last page (i.e., in the middle of the history), delete all forward pages
	if bh.current != nil && bh.current.next != nil {
		// Cut off all forward pages
		bh.current.next = nil
		bh.tail = bh.current
	}

	// Add the new page to the end of the history
	if bh.head == nil { // If the history is empty
		bh.head = newNode
		bh.tail = newNode
		bh.current = newNode
	} else { // If there's already history
		bh.tail.next = newNode
		newNode.prev = bh.tail
		bh.tail = newNode
		bh.current = newNode
	}
}

// Back navigates to the previous page in the history
func (bh *BrowserHistory) Back() error {
	if bh.current == nil || bh.current.prev == nil {
		return errors.New("no previous page")
	}
	bh.current = bh.current.prev
	return nil
}

// Forward navigates to the next page in the history
func (bh *BrowserHistory) Forward() error {
	if bh.current == nil || bh.current.next == nil {
		return errors.New("no next page")
	}
	bh.current = bh.current.next
	return nil
}

// ClearHistory clears the entire browser history
func (bh *BrowserHistory) ClearHistory() {
	bh.current = nil
	bh.head = nil
	bh.tail = nil
}

// GetCurrentURL returns the URL of the current page
func (bh *BrowserHistory) GetCurrentURL() string {
	if bh.current == nil {
		return ""
	}
	return bh.current.url
}

// Example usage
func main() {
	bh := NewBrowserHistory()
	bh.VisitNewPage("google.com")
	bh.VisitNewPage("facebook.com")
	bh.VisitNewPage("twitter.com")

	fmt.Println("Current URL:", bh.GetCurrentURL()) // twitter.com
	bh.Back()
	fmt.Println("Current URL after Back:", bh.GetCurrentURL()) // facebook.com
	bh.VisitNewPage("linkedin.com")
	fmt.Println("Current URL after visiting new page:", bh.GetCurrentURL()) // linkedin.com

	// Test: Forward pages should be deleted
	err := bh.Forward()
	if err != nil {
		fmt.Println("Error:", err) // should print no next
	}

	bh.ClearHistory()
	fmt.Println("Current URL after clearing history:", bh.GetCurrentURL()) // empty string
}
