package util

import (
	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/client"
	"github.com/theRealAlpaca/go-selenium/client/session"
	"github.com/theRealAlpaca/go-selenium/config"
	"github.com/theRealAlpaca/go-selenium/logger"
)

func HandleError(s *session.Session, err error) {
	if err == nil {
		return
	}

	if !errors.As(err, api.FailedRequestErr) {
		panic(err)
	}

	logger.Error(err)

	if config.Config.SoftAsserts {
		s.Errors = append(s.Errors, err.Error())

		return
	}

	logger.Error(err)

	if err := client.DeleteSession(s); err != nil {
		panic(errors.Wrap(err, "failed to delete session"))
	}

	client.Stop()

}
func HandleResponseError(s *session.Session, res *api.ErrorResponse) {
	if res == nil {
		return
	}

	if config.Config.SoftAsserts {
		s.Errors = append(s.Errors, res.String())

		return
	}

	logger.Error(res.String())

	if err := client.DeleteSession(s); err != nil {
		panic(errors.Wrap(err, "failed to delete session"))
	}

	// TODO: Maybe Stop only the current session and not the whole client?
	client.Stop()
}
