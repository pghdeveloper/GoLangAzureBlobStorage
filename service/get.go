package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

type GetFilesRepository interface {
	GetFilesFromCloud(ctx context.Context, containerId string) []string
}

var Repos GetFilesRepository

func GetFileNames(c *gin.Context) {
	containerId := c.Param("containerId")
	ctx := context.Background()

	strArray := Repos.GetFilesFromCloud(ctx, containerId)

	fmt.Println("Before checking Length of str array")
	fmt.Println("strArray: ", strArray)
	if len(strArray) == 0 {
		log.Fatal("Files not exist")
	}

	c.IndentedJSON(http.StatusOK, strArray)
}