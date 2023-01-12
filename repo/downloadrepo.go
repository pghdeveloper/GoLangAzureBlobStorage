package repo

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type AzureDownloadRepo struct {

}

func (Az *AzureDownloadRepo) DownloadFileFromCloud(ctx context.Context, containerId string, fileName string) (*bytes.Buffer, error) {
	
	_, accountPath, credential := ConnectNew()

	blobClient, err := azblob.NewBlockBlobClientWithSharedKey(accountPath+containerId+"/"+fileName, credential, nil)
	if err != nil {
		return new(bytes.Buffer), err
	}

	fmt.Println(accountPath+containerId+"/"+fileName)

	// Download the blob
	get, err := blobClient.Download(ctx, nil)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return new(bytes.Buffer), err
	}

	downloadedData := &bytes.Buffer{}
	reader := get.Body(azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(reader)
	if err != nil {
		return new(bytes.Buffer), err
	}
	err = reader.Close()
	if err != nil {
		return new(bytes.Buffer), err
	}

	return downloadedData, nil
}