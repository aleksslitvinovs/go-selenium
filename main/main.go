package main

import (
	"time"

	"github.com/theRealAlpaca/go-selenium"
	"github.com/theRealAlpaca/go-selenium/client"
	"github.com/theRealAlpaca/go-selenium/client/session/element"
	"github.com/theRealAlpaca/go-selenium/driver"
)

//nolint:errcheck
func main() {
	d := driver.NewDriver(
		"/Users/aleksslitvinovs/Downloads/chromedriver",
		4444,
		"http://localhost",
		&driver.Opts{
			Timeout: time.Second * 10,
		},
	)

	s := selenium.Start(d, nil)

	defer client.Stop()

	s.OpenURL("https://app.stage.loadero.com/login")

	element.
		NewElement(s, ".sign-in-form").
		WaitFor(time.Second * 5).
		UntilIsVisible()

	element.
		NewElement(s, "#username").
		WaitFor(time.Second * 5).
		UntilIsVisible().
		SendKeys("testing@loadero.abc")

	element.
		NewElement(s, "#password").
		WaitFor(time.Second * 5).
		UntilIsVisible().
		SendKeys("password")

	element.
		NewElement(s, ".button--primary").
		WaitFor(time.Second * 5).
		UntilIsVisible().
		Click()

	time.Sleep(10 * time.Second)
}
