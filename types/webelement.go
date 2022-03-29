package types

import "time"

type Waiterer interface {
	UntilIsVisible() WebElementer
	UntilIsNotVisible() WebElementer
	UntilIsEnabled() WebElementer
	UntilIsNotEnabled() WebElementer
	UntilIsSelected() WebElementer
	UntilIsNotSelected() WebElementer
}

type WebElementer interface {
	FindElement() string
	WaitFor(timeout time.Duration) Waiterer
	// TODO: Handle error
	GetText() (string, error)
	Click() error
	SendKeys(input string)
	Clear() error
	IsVisible() bool
	IsEnabled() bool
	IsSelected() bool
}
