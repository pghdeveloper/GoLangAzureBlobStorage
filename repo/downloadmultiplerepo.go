package repo

import (
	"bytes"
	"context"
	"example/GoLangAzureBlobStorage/lib"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type AzureDownloadMultipleRepo struct {

}

func (Az *AzureDownloadMultipleRepo) DownloadMultipleFilesFromCloud(ctx context.Context, containerIds lib.Containers) ([]*lib.InMemoryFile, error) {
	serviceClient, accountPath, credential := ConnectNew()

	inMemoryFiles := []*lib.InMemoryFile{}
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

			inMemoryFile := new(lib.InMemoryFile)
			inMemoryFile.Content = downloadedData.Bytes()
			inMemoryFile.FileName = blob

			inMemoryFiles = append(inMemoryFiles, inMemoryFile)
		}
	}

	return inMemoryFiles, nil
}