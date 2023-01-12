package service

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

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

func DownloadMultiple(c *gin.Context) {
	var containerIds Containers
	ctx := context.Background()

	if err := c.BindJSON(&containerIds); err != nil {
		return
	}

	serviceClient, accountPath, credential := Connect()

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

type DownloadRepository interface {
	DownloadFileFromCloud(ctx context.Context, containerId string, fileName string) (*bytes.Buffer, error)
}

var DownloadRepos DownloadRepository

func DownloadFile(c *gin.Context) {
	containerId := c.Param("containerId")
	fileName := c.Param("fileName")
	ctx := context.Background()

	fmt.Println("HI HI")
	downloadedData, err := DownloadRepos.DownloadFileFromCloud(ctx, containerId, fileName)
	if (err != nil) {
		log.Println("Error: " + err.Error())
		c.JSON(http.StatusNotFound, gin.H {
		 	"Message": "Issue Downloading the file",
		})
		return
	}

	fmt.Println(downloadedData.String())
	c.JSON(200, downloadedData.Bytes())
	//c.Header("Content-Disposition", "attachment; filename="+fileName)
	//c.Data(http.StatusOK, "application/octet-stream", downloadedData.Bytes())
}