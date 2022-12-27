package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetFileNames(c *gin.Context) {
	containerId := c.Param("containerId")
	ctx := context.Background()

	serviceClient, _, _ := Connect()

	containerClient := serviceClient.NewContainerClient(containerId)

	pager:= containerClient.ListBlobsFlat(nil)

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

	c.IndentedJSON(http.StatusOK, strArray)
}