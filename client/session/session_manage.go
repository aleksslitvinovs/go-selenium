package session

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/util"
)

func CreateSession(c api.Requester) string {
	req := struct {
		Capabilities map[string]interface{} `json:"capabilities"`
	}{
		make(map[string]interface{}),
	}

	var r struct {
		Value struct {
			SessionID    string                 `json:"sessionId"`
			Capabilities map[string]interface{} `json:"capabilities"`
		} `json:"value"`
	}

	err := api.ExecuteRequestCustom(http.MethodPost, "/session", c, req, &r)
	if err != nil {
		panic(errors.Wrap(err, "failed to start session"))
	}

	return r.Value.SessionID
}

func (s *Session) DeleteSession() {
	res, err := api.ExecuteRequestVoid(
		http.MethodDelete,
		fmt.Sprintf("/session/%s", s.ID),
		s,
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return
		}

		util.HandleError(s, errors.Wrap(err, "failed to delete session"))
	}

	s.KillDriver <- struct{}{}
}
