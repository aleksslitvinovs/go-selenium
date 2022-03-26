package main

import (
	"fmt"
	"time"

	"github.com/theRealAlpaca/go-selenium"
	"github.com/theRealAlpaca/go-selenium/client"
	"github.com/theRealAlpaca/go-selenium/client/session/element"
	"github.com/theRealAlpaca/go-selenium/client/session/element/selectors"
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

	err := s.OpenURL("https://duckduckgo.com/")
	if err != nil {
		panic(err)
	}

	url, err := s.GetCurrentURL()
	if err != nil {
		panic(err)
	}

	fmt.Println(url)

	element.NewElement(s, selectors.CSS, "test").FindElement()

	bodyElement := element.
		NewElement(s, selectors.CSS, "body").
		WaitFor(time.Second * 5).
		UntilIsVisible()

	err = bodyElement.SendKeys(s, "Hello World")
	if err != nil {
		panic(err)
	}

	clickButton := element.NewElement(s, selectors.CSS, "[type=submit]")
	clickButton.WaitFor(time.Second * 5).UntilIsVisible().Click(s)

	searchResultElement := element.NewElement(
		s, selectors.CSS, "#r1-0 .result__title",
	)
	searchResultElement.WaitFor(time.Second * 5).UntilIsVisible()
}
