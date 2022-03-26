package util

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/client"
	"github.com/theRealAlpaca/go-selenium/client/session"
)

func HandleError(s *session.Session, err error) {
	if s.Config.SoftAsserts {
		s.Errors = append(s.Errors, err)

		return
	}

	if err := client.DeleteSession(s); err != nil {
		panic(errors.Wrap(err, "failed to delete session"))
	}

	fmt.Println(err.Error())
}
