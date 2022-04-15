package selenium

import (
	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/logger"
	"github.com/theRealAlpaca/go-selenium/types"
)

func HandleError(err error) {
	if err == nil {
		return
	}

	logger.Error(err)

	if !errors.As(err, &types.ErrFailedRequest) {
		panic(err)
	}

	if Config.SoftAsserts {
		return
	}

	panic(err)
}

func HandleResponseError(res *ErrorResponse) {
	if res == nil {
		return
	}

	logger.Error(res.String())

	if Config.SoftAsserts {
		return
	}

	panic(res.String())
}
