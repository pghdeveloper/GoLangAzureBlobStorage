package repo

import (
	"example/GoLangAzureBlobStorage/util"
	"fmt"
	"log"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func ConnectNew() (azblob.ServiceClient, string, *azblob.SharedKeyCredential) {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	credential, err := azblob.NewSharedKeyCredential(config.AccountName, config.AccountKey)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}

	accountPath := fmt.Sprintf("https://%s.blob.core.windows.net/", config.AccountName)
	serviceClient, err := azblob.NewServiceClientWithSharedKey(accountPath, credential, nil)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}

	return serviceClient, accountPath, credential
}