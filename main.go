package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gin-gonic/gin"
)

type Container struct {
	ContainerId string
}

func downloadFile(c *gin.Context) {
	containerName := c.Param("containerId")
	fileName := c.Param("fileName")
	ctx := context.Background()
	credential, err := azblob.NewSharedKeyCredential("", "")
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}

	accountPath := fmt.Sprintf("https://%s.blob.core.windows.net/", "")

	blobClient, err := azblob.NewBlockBlobClientWithSharedKey(accountPath+containerName+"/"+fileName, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Download the blob
	get, err := blobClient.Download(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	downloadedData := &bytes.Buffer{}
	reader := get.Body(azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(reader)
	if err != nil {
		log.Fatal(err)
	}
	err = reader.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(downloadedData.String())
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Data(http.StatusOK, "application/pdf", downloadedData.Bytes())
}

func getFileNames(c *gin.Context) {
	containerId := c.Param("containerId")
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

	containerClient := serviceClient.NewContainerClient(containerId)

	pager := containerClient.ListBlobsFlat(nil)

	var strArray []string
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, v := range resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems {
			fmt.Println(*v.Name)
			strArray = append(strArray, *v.Name)
		}
	}

	if err = pager.Err(); err != nil {
		log.Fatalf("Failure to list blobs: %+v", err)
	}
	c.IndentedJSON(http.StatusOK, strArray)
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

	credential, err := azblob.NewSharedKeyCredential("", "")
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}

	accountPath := fmt.Sprintf("https://%s.blob.core.windows.net/", "")
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
		go send(fileHeader, accountPath, containerName, credential)
	}
	c.IndentedJSON(http.StatusOK, "Y")
}

func main() {
	router := gin.Default()
	router.GET("/getListOfDocumentsById/:containerId", getFileNames)
	router.POST("/uploadMultiple", sendToAzureFiles)
	router.GET("download/:containerId/:fileName", downloadFile)
	router.Run("localhost:8081")
}
