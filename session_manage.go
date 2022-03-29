package selenium

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/types"
	"github.com/theRealAlpaca/go-selenium/util"
)

func (c *client) CreateSession() (types.Sessioner, error) {
	err := c.waitUntilIsReady(10 * time.Second)
	if err != nil {
		return &session{}, errors.Wrap(
			err, "driver is not ready to start a new session",
		)
	}

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

	err = c.api.ExecuteRequestCustom(http.MethodPost, "/session", req, &r)
	if err != nil {
		panic(errors.Wrap(err, "failed to start session"))
	}

	s := &session{
		url: c.api.BaseURL,
		id:  r.Value.SessionID,
		api: &api.APIClient{BaseURL: c.api.BaseURL},
	}

	s.killDriver = make(chan struct{})

	c.sessions[s] = true

	go c.sessionListener(s)

	return s, nil
}

func (s *session) DeleteSession() {
	res, err := s.api.ExecuteRequestVoid(
		http.MethodDelete,
		fmt.Sprintf("/session/%s", s.id),
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			util.HandleResponseError(s, errRes)

			return
		}

		util.HandleError(s, errors.Wrap(err, "failed to delete session"))
	}

	s.killDriver <- struct{}{}
}

func (c *client) sessionListener(s *session) {
	<-s.killDriver

	delete(c.sessions, s)

	c.MustStop()
}
