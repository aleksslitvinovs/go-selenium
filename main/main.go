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
	)
	c := client.NewClient(d)

	s, err := selenium.Start(c)
	if err != nil {
		panic(err)
	}

	defer func() {
		// TODO: Kill GC after connection is closed
		err = c.Stop()
		if err != nil {
			panic(err)
		}
	}()

	err = s.OpenURL("https://duckduckgo.com/")
	if err != nil {
		panic(err)
	}

	url, err := s.GetCurrentURL()
	if err != nil {
		panic(err)
	}

	fmt.Println(url)

	bodyElement := element.NewElement(
		selectors.CSS, "#search_form_input_homepage",
	)

	err = bodyElement.SendKeys(s, "Hello World")
	if err != nil {
		panic(err)
	}

	clickButton := element.NewElement(selectors.CSS, "[type=submit]")
	clickButton.WaitFor(s, time.Second*5)
	clickButton.Click(s)

	searchResultElement := element.NewElement(
		selectors.CSS, "#r1-0 .result__title",
	)
	searchResultElement.WaitFor(s, time.Second*5).UntilIsVisible()

	text, err := searchResultElement.GetText(s)
	if err != nil {
		panic(err)
	}

	fmt.Println(text)

	s.Refresh()

	if err != nil {
		panic(err)
	}
}
