package selenium

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/selector"
)

type Session struct {
	id             string
	defaultLocator string
	// TODO: Maybe create a custom struct for handling error types. Maybe just
	// an alias to string? Maybe could implement Error interface?
	errors []string
	api    *APIClient
}

func NewSession() (*Session, error) {
	if Client == nil {
		return nil, errors.New("client is not set")
	}

	err := Client.waitUntilIsReady(10 * time.Second)
	if err != nil {
		return nil, errors.Wrap(
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
		return nil, errors.Wrap(err, "failed to start session")
	}

	s := &Session{
		id:             response.Value.SessionID,
		defaultLocator: Config.ElementSettings.SelectorType,
		api:            Client.api,
	}

	Client.sessions[s] = true

	return s, nil
}

func (s *Session) DeleteSession() {
	res, err := s.api.ExecuteRequestVoid(
		http.MethodDelete,
		fmt.Sprintf("/session/%s", s.id),
	)
	if err != nil {
		if errRes := res.GetErrorReponse(); errRes != nil {
			HandleResponseError(errRes)
		}

		HandleError(errors.Wrap(err, "failed to delete session"))
	}
}

func (s *Session) GetID() string {
	return s.id
}

func (s *Session) AddError(err string) {
	s.errors = append(s.errors, err)
}

func (s *Session) UseCSS() {
	s.defaultLocator = selector.CSS
}

func (s *Session) UseXPath() {
	s.defaultLocator = selector.XPath
}

func (s *Session) RaiseErrors() string {
	if len(s.errors) == 0 {
		return ""
	}

	errors := make([]string, 0, len(s.errors))

	errors = append(errors, s.errors...)

	return strings.Join(errors, "\n")
}

func getCapabilities() map[string]interface{} {
	caps := Config.WebDriver.Capabalities
	if caps == nil {
		caps = make(map[string]interface{})
	}

	finalCaps := make(map[string]interface{})

	finalCaps["alwaysMatch"] = caps

	return finalCaps
}
