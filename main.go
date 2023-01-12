package main

import (
	"example/GoLangAzureBlobStorage/repo"
	"example/GoLangAzureBlobStorage/service"
	"github.com/gin-gonic/gin"
)



func main() {
	service.Repos = &repo.AzureRepo {}
	service.DownloadRepos = &repo.AzureDownloadRepo {}

	router := gin.Default()
	router.GET("/getListOfDocumentsById/:containerId", service.GetFileNames)
	router.POST("/uploadMultiple", service.SendToAzureFiles)
	router.GET("download/:containerId/:fileName", service.DownloadFile)
	router.POST("downloadmultiple", service.DownloadMultiple)
	router.Run("localhost:8081")
}
