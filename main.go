package main

import (
	"bytes"
	"context"
	"example/GoLangAzureBlobStorage/service"
	"fmt"
	"log"
	"net/http"
	"archive/zip"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gin-gonic/gin"
)

type Containers struct {
	ContainerIds []string
}

type InMemoryFile struct {
	FileName string
	Content  []byte
}

func createZipFile(inMemoryFiles []*InMemoryFile) []byte{
	fmt.Println("we are in the zipData function")
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	zipWriter := zip.NewWriter(buf)

	for _, file := range inMemoryFiles {
		zipFile, err := zipWriter.Create(file.FileName)
		if err != nil {
			fmt.Println(err)
		}
		_, err = zipFile.Write(file.Content)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Make sure to check the error on Close.
	err := zipWriter.Close()
	if err != nil {
		fmt.Println(err)
	}

	//write the zipped file to the disk
	return buf.Bytes()
}

func downloadMultiple(c *gin.Context) {
	var containerIds Containers
	ctx := context.Background()

	if err := c.BindJSON(&containerIds); err != nil {
		return
	}

	serviceClient, accountPath, credential := service.Connect()

	inMemoryFiles := []*InMemoryFile{}
	for _, containerId := range containerIds.ContainerIds {
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

		if pager.Err() != nil {
			log.Fatalf("Failure to list blobs: %+v", pager.Err())
		}
		

		for _, blob := range strArray {
			fmt.Println(accountPath+containerId+"/"+blob)
			blobClient, err := azblob.NewBlockBlobClientWithSharedKey(accountPath+containerId+"/"+blob, credential, nil)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("About to Download")

			get, err := blobClient.Download(ctx, nil)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Done Downloading")

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

			inMemoryFile := new(InMemoryFile)
			inMemoryFile.Content = downloadedData.Bytes()
			inMemoryFile.FileName = blob

			inMemoryFiles = append(inMemoryFiles, inMemoryFile)
		}
	}

	zipFile := createZipFile(inMemoryFiles)

	c.JSON(200, zipFile)
	//c.Header("Content-Disposition", "attachment; filename=zipFile.zip")
    //c.Data(http.StatusOK, "application/octet-stream", zipFile)
}

func downloadFile(c *gin.Context) {
	containerName := c.Param("containerId")
	fileName := c.Param("fileName")
	ctx := context.Background()

	_, accountPath, credential := service.Connect()

	blobClient, err := azblob.NewBlockBlobClientWithSharedKey(accountPath+containerName+"/"+fileName, credential, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(accountPath+containerName+"/"+fileName)

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
	c.JSON(200, downloadedData.Bytes())
	//c.Header("Content-Disposition", "attachment; filename="+fileName)
	//c.Data(http.StatusOK, "application/octet-stream", downloadedData.Bytes())
}

func main() {
	router := gin.Default()
	router.GET("/getListOfDocumentsById/:containerId", service.GetFileNames)
	router.POST("/uploadMultiple", service.SendToAzureFiles)
	router.GET("download/:containerId/:fileName", downloadFile)
	router.POST("downloadmultiple", downloadMultiple)
	router.Run("localhost:8081")
}
