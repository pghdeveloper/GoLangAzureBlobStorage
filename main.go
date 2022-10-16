package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

type Container struct {
	ContainerId string
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

//func getBooks(c *gin.Context) {
//sent, err := sendToAzure()
//if err != nil {
//return
//}
//c.IndentedJSON(http.StatusOK, sent)
//}

func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
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

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available."})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func sendToAzureFiles(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		log.Fatal("File too large")
	}

	files := c.Request.MultipartForm.File["attachments"]
	//Parse Json String
	value := c.Request.FormValue("data")

	fmt.Println(value)

	// defining a struct instance
	var container1 Container

	// data in JSON format which
	// is to be decoded
	Data := []byte(value)

	// decoding container1 struct
	// from json format
	err := json.Unmarshal(Data, &container1)

	if err != nil {

		// if error is not nil
		// print error
		fmt.Println(err)
	}

	ctx := context.Background()

	credential, err := azblob.NewSharedKeyCredential("golangdocumentapisa", "CHZDQAHuHeEwUjhEkbpIWT0awgo/HI5WPhn7zYQVKKYcfXfESVlDW1uQmM6CMvmzTe/F7pEbXRsR+AStVzkJ/w==")
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}

	accountPath := fmt.Sprintf("https://%s.blob.core.windows.net/", "golangdocumentapisa")
	serviceClient, err := azblob.NewServiceClientWithSharedKey(accountPath, credential, nil)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}

	containerName := "golangcontainer" + "-" + container1.ContainerId

	containerClient := serviceClient.NewContainerClient(containerName)
	fmt.Println("HI2.5")
	_, err = containerClient.Create(ctx, nil)
	fmt.Println("HI2.6")
	if err != nil {
		fmt.Println("HI-Error")
		log.Fatal(err)
	}
	for _, fileHeader := range files {
		if fileHeader.Size > 1000000000 {
			log.Fatal("File Too Large")
		}

		// Open the file
		file, _ := fileHeader.Open()
		dat, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal("Problem Opening File")
		}

		fmt.Println(fileHeader.Filename)

		fmt.Println("HI")

		fmt.Println("HI2")

		fmt.Println("HI2.7")
		blobName := fileHeader.Filename

		fmt.Println("HI3")
		blobClient, err := azblob.NewBlockBlobClientWithSharedKey(accountPath+containerName+"/"+blobName, credential, nil)
		if err != nil {
			log.Fatal(err)
		}

		// Upload to data to blob storage
		_, err = blobClient.UploadBufferToBlockBlob(ctx, dat, azblob.HighLevelUploadToBlockBlobOption{})

		fmt.Println("HI4")
		if err != nil {
			log.Fatalf("Failure to upload to blob: %+v", err)
		}
		defer file.Close()
	}
}

func sendToAzure(c *gin.Context) {
	fmt.Println("Intro")
	file, _ := c.FormFile("file")

	//Parse Json String
	value := c.Request.FormValue("data")

	fmt.Println(value)

	// defining a struct instance
	var container1 Container

	// data in JSON format which
	// is to be decoded
	Data := []byte(value)

	// decoding container1 struct
	// from json format
	err := json.Unmarshal(Data, &container1)

	if err != nil {

		// if error is not nil
		// print error
		fmt.Println(err)
	}

	fileContent, _ := file.Open()
	dat, err := ioutil.ReadAll(fileContent)
	if err != nil {
		log.Fatal("Cannot read file " + err.Error())
	}

	fmt.Println(file.Filename)

	fmt.Println("HI")
	ctx := context.Background()

	credential, err := azblob.NewSharedKeyCredential("", "")
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}

	accountPath := fmt.Sprintf("https://%s.blob.core.windows.net/", "")
	serviceClient, err := azblob.NewServiceClientWithSharedKey(accountPath, credential, nil)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}

	fmt.Println("HI2")
	containerName := "golangcontainer" + "-" + container1.ContainerId
	containerClient := serviceClient.NewContainerClient(containerName)

	fmt.Println("HI2.5")
	_, err = containerClient.Create(ctx, nil)
	fmt.Println("HI2.6")
	if err != nil {
		fmt.Println("HI-Error")
		log.Fatal(err)
	}

	fmt.Println("HI2.7")
	blobName := file.Filename

	fmt.Println("HI3")
	blobClient, err := azblob.NewBlockBlobClientWithSharedKey(accountPath+containerName+"/"+blobName, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Upload to data to blob storage
	_, err = blobClient.UploadBufferToBlockBlob(ctx, dat, azblob.HighLevelUploadToBlockBlobOption{})

	fmt.Println("HI4")
	if err != nil {
		log.Fatalf("Failure to upload to blob: %+v", err)
	}

	c.IndentedJSON(http.StatusOK, "Y")
}

func main() {
	router := gin.Default()
	//router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.POST("/upload", sendToAzure)
	router.POST("/uploadMultiple", sendToAzureFiles)
	router.Run("localhost:8081")
}
