package elements

import (
	"github.com/theRealAlpaca/go-selenium/client"
	"github.com/theRealAlpaca/go-selenium/client/element"
)

type Elements struct {
	SelectorType string `json:"using"`
	Selector     string `json:"value"`

	Values []*element.Element `json:"value"` //nolint:govet
}

func NewElements(selectorType, selector string) *Elements {
	return &Elements{
		SelectorType: selectorType,
		Selector:     selector,
	}
}

// TODO: Implement.
func (ee *Elements) FindElements(c *client.Client) ([]string, error) {
	return []string{}, nil
}
