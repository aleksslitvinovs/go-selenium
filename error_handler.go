package selenium

import (
	"github.com/theRealAlpaca/go-selenium/logger"
)

func handleError(res *response, err error) {
	if res == nil {
		logger.Error(err.Error())

		if config.SoftAsserts {
			return
		}

		panic(err.Error())
	}

	errRes := res.getErrorReponse()
	if errRes != nil {
		logger.Error(errRes)

		if config.SoftAsserts {
			return
		}

		// TODO: Handle error response (element not visible, etc.)
		panic(errRes.Error())
	}

	panic(err.Error())
}
