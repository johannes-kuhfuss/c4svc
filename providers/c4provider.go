package providers

import (
	"bufio"
	"os"

	"github.com/johannes-kuhfuss/c4/utils/logger"
	rest_errors "github.com/johannes-kuhfuss/c4/utils/rest_errors_utils"
)

var (
	C4Provider c4ProviderInterface = &c4ProviderService{}
)

type c4ProviderService struct{}

type c4ProviderInterface interface {
	ProcessFile(string, string, bool) (string, rest_errors.RestErr)
}

func (c4p *c4ProviderService) ProcessFile(srcUrl string, dstUrl string, rename bool) (string, rest_errors.RestErr) {
	file, openErr := os.Open(srcUrl)
	if openErr != nil {
		logger.Error("Source file not found or could not be read", openErr)
		return "", rest_errors.NewNotFoundError("Source file not found or could not be read")
	}
	defer file.Close()
	fileReader := bufio.NewReader(file)
	_ = fileReader
	/*
		if id == nil {
			logger.Error(fmt.Sprintf("Could not generate C4 Id for %v", srcUrl), nil)
			return "", rest_errors.NewInternalServerError("could not generate C4 Id", nil)
		}
	*/
	return "", nil
}
