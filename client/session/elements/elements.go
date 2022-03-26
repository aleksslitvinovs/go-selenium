package elements

import (
	"github.com/theRealAlpaca/go-selenium/client/session/element"
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
