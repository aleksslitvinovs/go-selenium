package element

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/client"
)

func (e *Element) IsVisible(c *client.Client) (bool, error) {
	if e.webID == "" {
		id, err := e.FindElement(c)
		if err != nil {
			return false, errors.Wrap(err, "could not find element")
		}

		e.webID = id
	}

	res, err := client.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/displayed", c.SessionID, e.webID),
		e,
		c.Driver,
	)
	if err != nil {
		return false, errors.Wrap(err, "could not get display state")
	}

	var response struct {
		Value bool `json:"value"`
	}

	err = json.Unmarshal(res, &response)
	if err != nil {
		return false, errors.Wrap(err, "could not unmarshal response")
	}

	fmt.Println("Visibility result " + string(res))

	return response.Value, nil
}

func (e *Element) IsEnabled(c *client.Client) (bool, error) {
	if e.webID == "" {
		id, err := e.FindElement(c)
		if err != nil {
			return false, errors.Wrap(err, "could not find element")
		}

		e.webID = id
	}

	res, err := client.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/enabled", c.SessionID, e.webID),
		e,
		c.Driver,
	)
	if err != nil {
		return false, errors.Wrap(err, "could not get enabled stated")
	}

	var response struct {
		Value bool `json:"value"`
	}

	err = json.Unmarshal(res, &response)
	if err != nil {
		return false, errors.Wrap(err, "could not unmarshal response")
	}

	fmt.Println("Enabled result " + string(res))

	return response.Value, nil
}

func (e *Element) IsSelected(c *client.Client) (bool, error) {
	if e.webID == "" {
		id, err := e.FindElement(c)
		if err != nil {
			return false, errors.Wrap(err, "could not find element")
		}

		e.webID = id
	}

	res, err := client.ExecuteRequest(
		http.MethodGet,
		fmt.Sprintf("/session/%s/element/%s/selected", c.SessionID, e.webID),
		e,
		c.Driver,
	)
	if err != nil {
		return false, errors.Wrap(err, "could not get selected state")
	}

	var response struct {
		Value bool `json:"value"`
	}

	err = json.Unmarshal(res, &response)
	if err != nil {
		return false, errors.Wrap(err, "could not unmarshal response")
	}

	fmt.Println("Selected result " + string(res))

	return response.Value, nil
}
