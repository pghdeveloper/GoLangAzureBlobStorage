package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gin-gonic/gin"
)

type Container struct {
	ContainerId string
}

func send(fileHeader *multipart.FileHeader, accountPath string, containerName string, credential *azblob.SharedKeyCredential) {
	ctx := context.Background()

	if fileHeader.Size > 1000000000 {
		log.Fatal("File Too Large")
	}

	// Open the file
	file, _ := fileHeader.Open()
	dat, err := io.ReadAll(file)
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
	amt := time.Duration(rand.Intn(250))
	time.Sleep(time.Millisecond * amt)
}

func SendToAzureFiles(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		log.Fatal("File too large")
	}

	files := c.Request.MultipartForm.File["attachments"]
	//Parse Json String
	value := c.Request.FormValue("data")

	fmt.Println(value)

	// defining a struct instance
	var container Container

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

	serviceClient, accountPath, credential := Connect()

	containerName := "golangcontainer" + "-" + container.ContainerId

	containerClient := serviceClient.NewContainerClient(containerName)
	fmt.Println("HI2.5")
	_, err = containerClient.Create(ctx, nil)
	fmt.Println("HI2.6")
	if err != nil {
		fmt.Println("HI-Error")
		if !strings.Contains(err.Error(), "ContainerAlreadyExists") {
			log.Fatal(err)
		}	
	}
	for _, fileHeader := range files {
		go send(fileHeader, accountPath, containerName, credential)
	}
	c.IndentedJSON(http.StatusOK, "Y")
}

