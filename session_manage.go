package selenium

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/config"
	"github.com/theRealAlpaca/go-selenium/types"
	"github.com/theRealAlpaca/go-selenium/util"
)

func CreateSession() (types.Sessioner, error) {
	if Client == nil {
		return nil, errors.New("client is not set")
	}

	err := Client.waitUntilIsReady(10 * time.Second)
	if err != nil {
		return &Session{}, errors.Wrap(
			err, "driver is not ready to start a new session",
		)
	}

	req := struct {
		Capabilities map[string]interface{} `json:"capabilities"`
	}{getCapabilities()}

	var response struct {
		Value struct {
			SessionID    string                 `json:"sessionId"`
			Capabilities map[string]interface{} `json:"capabilities"`
		} `json:"value"`
	}

	err = Client.api.ExecuteRequestCustom(
		http.MethodPost, "/session", req, &response,
	)
	if err != nil {
		panic(errors.Wrap(err, "failed to start session"))
	}

	s := &Session{
		url: Client.api.BaseURL,
		id:  response.Value.SessionID,
		api: &api.APIClient{BaseURL: Client.api.BaseURL},
	}

	s.killDriver = make(chan struct{})

	Client.sessions[s] = true

	Client.runner.session = s

	go Client.sessionListener(s)

	return s, nil
}

func (s *Session) DeleteSession() {
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

func getCapabilities() map[string]interface{} {
	caps := config.Config.WebDriver.Capabalities
	if caps == nil {
		caps = make(map[string]interface{})
	}

	finalCaps := make(map[string]interface{})

	finalCaps["alwaysMatch"] = caps

	return finalCaps
}

func (c *client) sessionListener(s *Session) {
	// TODO: Kill only the session that is being used, not the whole driver.
	<-s.killDriver

	delete(c.sessions, s)

	MustStopClient()
}
