package util

import (
	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/config"
	"github.com/theRealAlpaca/go-selenium/logger"
	"github.com/theRealAlpaca/go-selenium/types"
)

func HandleError(s types.Sessioner, err error) {
	if err == nil {
		return
	}

	logger.Error(err)

	if !errors.As(err, &types.ErrFailedRequest) {
		panic(err)
	}

	if config.Config.SoftAsserts {
		return
	}

	panic(err)
}

func HandleResponseError(s types.Sessioner, res *api.ErrorResponse) {
	if res == nil {
		return
	}

	logger.Error(res.String())

	if config.Config.SoftAsserts {
		return
	}

	panic(res.String())
}
