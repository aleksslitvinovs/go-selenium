package element

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/client"
)

type Element struct {
	SelectorType string `json:"using"`
	Selector     string `json:"value"`
	webID        string `json:"-"`
}

func NewElement(selectorType, selector string) *Element {
	return &Element{
		SelectorType: selectorType,
		Selector:     selector,
	}
}

func (e *Element) GetText(c *client.Client) (string, error) {
	if e.webID == "" {
		id, err := e.FindElement(c)
		if err != nil {
			return "", errors.Wrap(err, "failed to get element")
		}

		e.webID = id
	}

	res, err := client.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/text", c.SessionID, e.webID),
		struct{}{},
		c.Driver,
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to send request to get text")
	}

	var response struct {
		Value string `json:"value"`
	}

	err = json.Unmarshal(res, &response)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal response")
	}

	fmt.Println("Get text response", string(res))

	return response.Value, nil
}

func (e *Element) Click(c *client.Client) error {
	if e.webID == "" {
		id, err := e.FindElement(c)
		if err != nil {
			return errors.Wrap(err, "failed to get element")
		}

		e.webID = id
	}

	_, err := client.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/click", c.SessionID, e.webID),
		e,
		c.Driver,
	)
	if err != nil {
		return errors.Wrap(err, "failed to send request to click")
	}

	return nil
}

func (e *Element) SendKeys(c *client.Client, input string) error {
	if e.webID == "" {
		id, err := e.FindElement(c)
		if err != nil {
			return errors.Wrap(err, "failed to get element")
		}

		e.webID = id
	}

	payload := struct {
		Text string `json:"text"`
	}{
		input,
	}

	_, err := client.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/value", c.SessionID, e.webID),
		payload,
		c.Driver,
	)
	if err != nil {
		return errors.Wrap(err, "failed to send request to send keys")
	}

	return nil
}

func (e *Element) Clear(c *client.Client) error {
	if e.webID == "" {
		id, err := e.FindElement(c)
		if err != nil {
			return errors.Wrap(err, "failed to get element")
		}

		e.webID = id
	}

	_, err := client.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element/%s/clear", c.SessionID, e.webID),
		e,
		c.Driver,
	)
	if err != nil {
		return errors.Wrap(err, "failed to send request to clear")
	}

	return nil
}
