package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Book Type.
type Book struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

// list of books data.
var books = []Book{
	{Id: "1", Title: "Title 1", Author: "Author 1", Quantity: 1},
	{Id: "2", Title: "Title 2", Author: "Author 2", Quantity: 2},
	{Id: "3", Title: "Title 3", Author: "Author 3", Quantity: 3},
}

func getBooks(c *gin.Context) { // gin.Context its same as request in django. it helps us also retunrn JSON.
	c.JSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var newBook Book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.JSON(http.StatusCreated, newBook)
}

func getBookById(id string) (*Book, error) {
	for i, book := range books {
		if book.Id == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("Book not found")
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, book)
}

func checkOutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Book id must exist.",
		})
		return
	}

	book, bookErr := getBookById(id)

	if bookErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": bookErr.Error(),
		})
		return
	}

	if book.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The given book is out of stock.",
		})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Book id must exist.",
		})
		return
	}

	book, bookErr := getBookById(id)

	if bookErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": bookErr.Error(),
		})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func main() {
	r := gin.Default()                       // initializes gin router.
	r.GET("/books", getBooks)                // GET /books router handler.
	r.POST("/books/create", createBook)      // POST /books/create router handler.
	r.GET("/books/:id", bookById)            // GET /books/:id router handler.
	r.PATCH("/books/checkout", checkOutBook) // PATCH /books/checkout?id=:id router handler.
	r.PATCH("/books/return", returnBook)     // PATCH /books/return?id=:id router handler.
	r.Run("localhost:8000")                  // by default it runs on port :8080.
}
