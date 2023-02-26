package service

import (
	"context"
	"encoding/json"
	"example/GoLangAzureBlobStorage/lib"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadRepository interface {
	UploadFilesToCloud(ctx context.Context, container lib.Container, files []*multipart.FileHeader)
}

var UploadRepos UploadRepository

func UploadFiles(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		log.Fatal("File too large")
	}

	files := c.Request.MultipartForm.File["attachments"]
	//Parse Json String
	value := c.Request.FormValue("data")

	fmt.Println(value)

	// defining a struct instance
	var container lib.Container

	// data in JSON format which
	// is to be decoded
	Data := []byte(value)

	// decoding container1 struct
	// from json format
	err := json.Unmarshal(Data, &container)

	if err != nil {

		// if error is not nil
		// print error
		fmt.Println(err)
	}

	ctx := context.Background()
	UploadRepos.UploadFilesToCloud(ctx, container, files)

	c.IndentedJSON(http.StatusOK, "Y")
}