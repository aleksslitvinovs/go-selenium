package util

import (
	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/config"
	"github.com/theRealAlpaca/go-selenium/logger"
)

type Sessioner interface {
	AddError(err string)
	DeleteSession()
}

func HandleError(s Sessioner, err error) {
	if err == nil {
		return
	}

	if !errors.As(err, &api.ErrFailedRequest) {
		panic(err)
	}

	s.AddError(err.Error())

	if config.Config.SoftAsserts {
		return
	}

	logger.Error(err)

	s.DeleteSession()
}
func HandleResponseError(s Sessioner, res *api.ErrorResponse) {
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
