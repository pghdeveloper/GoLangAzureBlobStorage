package service

import (
	"archive/zip"
	"bytes"
	"context"
	"example/GoLangAzureBlobStorage/lib"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DownloadRepository interface {
	DownloadFileFromCloud(ctx context.Context, containerId string, fileName string) (*bytes.Buffer, error)
}

type DownloadMultipleRepository interface {
	DownloadMultipleFilesFromCloud(ctx context.Context, containerIds lib.Containers) ([]*lib.InMemoryFile, error)
}

var DownloadRepos DownloadRepository
var DownloadMultipleRepos DownloadMultipleRepository

func createZipFile(inMemoryFiles []*lib.InMemoryFile) []byte{
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

func DownloadMultiple(c *gin.Context) {
	var containerIds lib.Containers
	ctx := context.Background()

	if err := c.BindJSON(&containerIds); err != nil {
		return
	}

	inMemoryFiles, err := DownloadMultipleRepos.DownloadMultipleFilesFromCloud(ctx, containerIds)
	if (err != nil) {
		log.Println("Error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H {
		 	"Message": "Issue Downloading file(s)",
		})
		return
	}

	zipFile := createZipFile(inMemoryFiles)

	c.JSON(200, zipFile)
	//c.Header("Content-Disposition", "attachment; filename=zipFile.zip")
    //c.Data(http.StatusOK, "application/octet-stream", zipFile)
}

func DownloadFile(c *gin.Context) {
	containerId := c.Param("containerId")
	fileName := c.Param("fileName")
	ctx := context.Background()

	fmt.Println("HI HI")
	downloadedData, err := DownloadRepos.DownloadFileFromCloud(ctx, containerId, fileName)
	if (err != nil) {
		log.Println("Error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H {
		 	"Message": "Issue Downloading the file",
		})
		return
	}

	fmt.Println(downloadedData.String())
	c.JSON(200, downloadedData.Bytes())
	//c.Header("Content-Disposition", "attachment; filename="+fileName)
	//c.Data(http.StatusOK, "application/octet-stream", downloadedData.Bytes())
}