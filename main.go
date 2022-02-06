package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Book Structure Serealized with JSON

type book struct {
	ID       string `json: "id"`
	Title    string `json: "title"`
	Author   string `json: "author"`
	Quantity int    `json: "quantity"`
}

var books = []book{ // Struct Slice
	{ID: "1", Title: "Hacking: The Art of Exploitation", Author: "Jon Erickson", Quantity: 2},
	{ID: "2", Title: "The Hash Crack: Password Cracking Manual", Author: "Joshua Picoletd", Quantity: 5},
	{ID: "3", Title: "The Hacker Playbook", Author: "Peter Kim", Quantity: 6},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)

}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not available"})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)

}

func returnBook(c *gin.Context) {

	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
	}

	book.Quantity += 1

	c.IndentedJSON(http.StatusOK, book)

}

func getBookById(id string) (*book, error) {
	for i, book := range books {
		if book.ID == id {
			return &books[i], nil // Having the pointer, allow us to modify the content of the book
		}
	}

	return nil, errors.New("book not found")
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, book)

}

func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)

}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")

}
