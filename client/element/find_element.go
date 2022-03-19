package element

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/client"
)

func (e *Element) FindElement(c *client.Client) (string, error) {
	res, err := client.ExecuteRequest(
		http.MethodPost,
		fmt.Sprintf("/session/%s/element", c.SessionID),
		e,
		c.Driver,
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to send request to get element")
	}

	var response struct {
		Value map[string]string `json:"value"`
	}

	err = json.Unmarshal(res, &response)
	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal response")
	}

	fmt.Println("Get element response", string(res))

	for _, v := range response.Value {
		if v != "" {
			return v, nil
		}
	}

	return "", errors.New("failed to get element id")
}
