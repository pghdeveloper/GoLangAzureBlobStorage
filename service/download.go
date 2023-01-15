package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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