package main

import(
	"github.com/gin-gonic/gin"
	"net/http"
	"errors"
	"strconv"
	"github.com/keploy/go-sdk/keploy"
	"github.com/keploy/go-sdk/integrations/kgin/v1"
)

type book struct{
	ID int `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Quantity int `json:"quantity"`
}

var books = []book{
	{ID: 1, Title: "The Alchemist", Author: "Paulo Coelho", Quantity: 10},
	{ID: 2, Title: "The Monk Who Sold His Ferrari", Author: "Robin Sharma", Quantity: 5},
	{ID: 3, Title: "The Power of Your Subconscious Mind", Author: "Joseph Murphy", Quantity: 7},
}

func addBook(c *gin.Context){
	var newBook book
	if err := c.ShouldBindJSON(&newBook); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	books = append(books, newBook)
	c.JSON(http.StatusCreated, newBook)
}

func getBooks(c *gin.Context){
	c.JSON(http.StatusOK, books)
}

func getBookById(id int)(*book, error){
	for i, b := range books{
		if(b.ID == id){
			return &books[i], nil
		}
	}
	return nil, errors.New("Book not found")
}

func bookById(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	book, err := getBookById(id)
	if err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, book)
}

func main(){

	k := keploy.New(keploy.Config{
	App: keploy.AppConfig{
		Name: "go-api-tutorial",
		Port: "8080",
	},
	Server: keploy.ServerConfig{
		URL: "http://localhost:6789/api",
	},

})
	router := gin.Default()
	kgin.GinV1(k, router)
	router.GET("/books", getBooks)
	router.POST("/book", addBook)
	router.GET("/book/:id", bookById)
	router.Run(":8080")
}

