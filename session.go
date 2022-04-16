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
	api    *apiClient
}

func NewSession() (*Session, error) {
	if client == nil {
		return nil, errors.New("client is not set")
	}

	err := client.waitUntilIsReady(10 * time.Second)
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

	res, err := client.api.executeRequestCustom(
		http.MethodPost, "/session", req, &response,
	)
	if err != nil {
		handleError(res, err)
	}

	s := &Session{
		id:             response.Value.SessionID,
		defaultLocator: config.ElementSettings.SelectorType,
		api:            client.api,
	}

	client.ss.mu.Lock()
	defer client.ss.mu.Unlock()

	client.ss.sessions[s] = true

	return s, nil
}

func (s *Session) DeleteSession() {
	res, err := s.api.executeRequestVoid(
		http.MethodDelete,
		fmt.Sprintf("/session/%s", s.id),
	)
	if err != nil {
		handleError(res, err)
	}

	client.ss.mu.Lock()
	defer client.ss.mu.Unlock()

	client.ss.sessions[s] = false
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
	caps := config.WebDriver.Capabalities
	if caps == nil {
		caps = make(map[string]interface{})
	}

	finalCaps := make(map[string]interface{})

	finalCaps["alwaysMatch"] = caps

	return finalCaps
}
