package webelement

import (
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/config"
	"github.com/theRealAlpaca/go-selenium/types"
)

type webElement struct {
	Selector     string `json:"value"`
	SelectorType string `json:"using"`

	id       string
	session  types.Sessioner
	settings *config.ElementSettings
	api      *api.APIClient
}

var _ (types.WebElementer) = (*webElement)(nil)

func NewElement(
	id string,
	session types.Sessioner,
	selector string,
	settings *config.ElementSettings,
	apiClient *api.APIClient,
) types.WebElementer {
	return &webElement{
		id:           id,
		session:      session,
		Selector:     selector,
		SelectorType: settings.SelectorType,
		settings:     settings,
		api:          &api.APIClient{BaseURL: apiClient.BaseURL},
	}
}
