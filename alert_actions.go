package selenium

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// DismissAlert dismisses currently open alert dialog.
func (s *Session) DismissAlert() *Session {
	res, err := s.api.executeRequestVoid(
		http.MethodPost, fmt.Sprintf("/session/%s/alert/dismiss", s.id),
	)
	if err != nil {
		handleError(res, err)

		return s
	}

	return s
}

// AcceptAlert accepts currently open alert dialog.
func (s *Session) AcceptAlert() *Session {
	res, err := s.api.executeRequestVoid(
		http.MethodPost, fmt.Sprintf("/session/%s/alert/accept", s.id),
	)
	if err != nil {
		handleError(res, err)

		return s
	}

	return s
}

// GetAlertText gets text of currently open alert dialog.
func (s *Session) GetAlertText() string {
	res, err := s.api.executeRequestVoid(
		http.MethodGet, fmt.Sprintf("/session/%s/alert/text", s.id),
	)
	if err != nil {
		handleError(res, err)

		return ""
	}

	if res.Value == nil {
		handleError(nil, errors.New("failed to get alert text"))
	}

	if v, ok := res.Value.(string); ok {
		return v
	}

	return ""
}

// SendAlertText sends text to currently open prompt dialog.
func (s *Session) SendAlertText(text string) *Session {
	payload := struct {
		Text string `json:"text"`
	}{
		Text: text,
	}

	res, err := s.api.executeRequest(
		http.MethodPost, fmt.Sprintf("/session/%s/alert/text", s.id), payload,
	)
	if err != nil {
		handleError(res, err)

		return s
	}

	return s
}
