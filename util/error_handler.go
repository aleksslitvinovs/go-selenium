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

	if !errors.As(err, &types.ErrFailedRequest) {
		panic(err)
	}

	s.AddError(err.Error())

	if config.Config.SoftAsserts {
		return
	}

	logger.Error(err)

	s.DeleteSession()
}

func HandleResponseError(s types.Sessioner, res *api.ErrorResponse) {
	if res == nil {
		return
	}

	s.AddError(res.String())

	if config.Config.SoftAsserts {
		return
	}

	logger.Error(res.String())

	s.DeleteSession()
}
