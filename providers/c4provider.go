package providers

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	c4gen "github.com/Avalanche-io/c4/id"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/johannes-kuhfuss/c4svc/config"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

var (
	C4Provider c4ProviderInterface = &c4ProviderService{}
)

type c4ProviderService struct{}

type c4ProviderInterface interface {
	ProcessFile(string, bool) (*string, *string, api_error.ApiErr)
}

func (c4p *c4ProviderService) ProcessFile(srcUrl string, rename bool) (*string, *string, api_error.ApiErr) {
	var dstUrl string
	if strings.TrimSpace(config.StorageAccountName) == "" || strings.TrimSpace(config.StorageAccountKey) == "" {
		logger.Error("No storage account access credentials", nil)
		return nil, nil, api_error.NewInternalServerError("No storage account access credentials", nil)
	}
	url, err := url.Parse(srcUrl)
	if err != nil || srcUrl == "" {
		logger.Error("Cannot parse source URL", nil)
		return nil, nil, api_error.NewBadRequestError("Cannot parse source URL")
	}
	blobUrl := url.Scheme + "://" + url.Host + "/"
	logger.Debug(fmt.Sprintf("blobUrl: %v", blobUrl))
	containerName := strings.TrimLeft(filepath.Dir(url.Path), string(os.PathSeparator))
	logger.Debug(fmt.Sprintf("containerName: %v", containerName))
	fileName := filepath.Base(url.Path)
	fileExt := filepath.Ext(url.Path)
	if url.Scheme == "" || url.Host == "" || containerName == "." {
		logger.Error("Cannot parse source URL", nil)
		return nil, nil, api_error.NewBadRequestError("Cannot parse source URL")
	}
	cred, err := azblob.NewSharedKeyCredential(config.StorageAccountName, config.StorageAccountKey)
	if err != nil {
		logger.Error("Cannot access storage account - wrong credentials", err)
		return nil, nil, api_error.NewInternalServerError("Cannot access storage account - wrong credentials", err)
	}

	//serviceClient, err := azblob.NewServiceClient(blobUrl, cred, nil)
	serviceClient, err := azblob.NewServiceClientWithSharedKey(blobUrl, cred, nil)

	if err != nil {
		logger.Error("Cannot access storage account - could not create service client", err)
		return nil, nil, api_error.NewInternalServerError("Cannot access storage account - could not create service client", err)
	}
	ctx := context.Background()
	container := serviceClient.NewContainerClient(containerName)
	blockBlob := container.NewBlobClient(fileName)
	get, err := blockBlob.Download(ctx, nil)
	if err != nil {
		logger.Error("Cannot access file on storage account", err)
		return nil, nil, api_error.NewBadRequestError("Cannot access file on storage account")
	}
	reader := get.Body(azblob.RetryReaderOptions{})
	id := c4gen.Identify(reader)
	c4string := id.String()
	if rename {
		lease, err := blockBlob.NewBlobLeaseClient(nil)
		if err != nil {
			logger.Error("Cannot get lease on file", err)
			return nil, nil, api_error.NewInternalServerError("Cannot get lease on file", err)
		}
		lease.AcquireLease(ctx, &azblob.AcquireLeaseBlobOptions{})
		defer lease.BreakLease(ctx, nil)
		newFileName := id.String() + fileExt
		newBlockBlob := container.NewBlobClient(newFileName)
		_, err = newBlockBlob.StartCopyFromURL(ctx, blockBlob.URL(), nil)
		if err != nil {
			logger.Error("Renaming of file failed", err)
			return nil, nil, api_error.NewInternalServerError("Renaming of file failed", err)
		}
		_, err = blockBlob.Delete(ctx, nil)
		if err != nil {
			logger.Error("Deleting of source file failed", err)
			return nil, nil, api_error.NewInternalServerError("Deleting of source file failed", err)
		}
		dstUrl = blobUrl + containerName + "/" + c4string + fileExt
	}
	if rename {
		return &c4string, &dstUrl, nil
	} else {
		return &c4string, nil, nil
	}
}
