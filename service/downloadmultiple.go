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

type DownloadMultipleRepository interface {
	DownloadMultipleFilesFromCloud(ctx context.Context, containerIds lib.Containers) ([]*lib.InMemoryFile, error)
}

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
	fmt.Println("Begin Process of Downloading Multiple")
	var containerIds lib.Containers
	ctx := context.Background()

	fmt.Println(c.Request.Body)
	fmt.Println("About to Bind Json")

	err := c.BindJSON(&containerIds)
	fmt.Println(containerIds)
	fmt.Println(len(containerIds.ContainerIds))

	if (len(containerIds.ContainerIds) == 0) {
		fmt.Println("HIHIHI")
	}
	fmt.Println("Before Error check")

	if (err != nil || len(containerIds.ContainerIds) == 0) {
		if (err != nil) {
			fmt.Println("Error Binding JSON: " + err.Error())
		}
		if (len(containerIds.ContainerIds) == 0) {
			fmt.Println("Error Binding JSON: Issue with Json Request")
		}
		c.JSON(http.StatusBadRequest, gin.H {
			"Message": "Json Request not correct",
		})
		return
	}

	fmt.Println(containerIds)
	fmt.Println("About to go to Download Multiple Repo")
	inMemoryFiles, err := DownloadMultipleRepos.DownloadMultipleFilesFromCloud(ctx, containerIds)
	fmt.Println("Done downloading from Download Multiple Repo")
	if (err != nil) {
		log.Println("Error: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H {
		 	"Message": "Issue Downloading file(s)",
		})
		return
	}

	zipFile := createZipFile(inMemoryFiles)
	fmt.Println("About to return Zip File to Client")

	c.JSON(200, zipFile)
	//c.Header("Content-Disposition", "attachment; filename=zipFile.zip")
    //c.Data(http.StatusOK, "application/octet-stream", zipFile)
}