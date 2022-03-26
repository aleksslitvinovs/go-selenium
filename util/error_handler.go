package util

import (
	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/config"
	"github.com/theRealAlpaca/go-selenium/logger"
)

type Halter interface {
	AddError(err string)
	Stop()
}

func HandleError(h Halter, err error) {
	if err == nil {
		return
	}

	if !errors.As(err, &api.ErrFailedRequest) {
		panic(err)
	}

	h.AddError(err.Error())

	if config.Config.SoftAsserts {
		return
	}

	logger.Error(err)
	// TODO: Maybe Stop only the current session and not the whole client?
	h.Stop()
}
func HandleResponseError(h Halter, res *api.ErrorResponse) {
	if res == nil {
		return
	}

	h.AddError(res.String())

	if config.Config.SoftAsserts {
		return
	}

	logger.Error(res.String())

	// TODO: Maybe Stop only the current session and not the whole client?
	h.Stop()
}
