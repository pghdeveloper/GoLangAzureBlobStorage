package repo

import (
	"context"
	"fmt"
)

type AzureRepo struct {

}

func (Az *AzureRepo) GetFilesFromCloud(ctx context.Context, containerId string) []string {
	
	serviceClient, _, _ := ConnectNew()

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

	return strArray
}