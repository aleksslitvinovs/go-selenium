package types

import "time"

type Waiterer interface {
	UntilIsPresent() WebElementer
	UntilIsNotPresent() WebElementer
	UntilIsVisible() WebElementer
	UntilIsNotVisible() WebElementer
	UntilIsEnabled() WebElementer
	UntilIsNotEnabled() WebElementer
	UntilIsSelected() WebElementer
	UntilIsNotSelected() WebElementer
}

type WebElementer interface {
	WaitFor(timeout time.Duration) Waiterer
	GetText() string
	Click() WebElementer
	SendKeys(input string) WebElementer
	Clear() WebElementer
	IsVisible() bool
	IsEnabled() bool
	IsSelected() bool
}
