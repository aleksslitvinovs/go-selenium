package main

import (
	"fmt"
	"time"

	"github.com/theRealAlpaca/go-selenium"
	"github.com/theRealAlpaca/go-selenium/client"
	"github.com/theRealAlpaca/go-selenium/client/element"
	"github.com/theRealAlpaca/go-selenium/client/element/selectors"
	"github.com/theRealAlpaca/go-selenium/driver"
)

func main() {
	d := driver.
		NewDriverBuilder().
		SetDriver("/Users/aleksslitvinovs/Downloads/chromedriver").
		SetPort(4444).
		Build()

	c := client.NewClientBuilder().
		SetDriver(d).
		Build()

	err := selenium.Start(c)
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

	err = c.OpenURL("https://duckduckgo.com/")
	if err != nil {
		panic(err)
	}

	url, err := c.GetURL()
	if err != nil {
		panic(err)
	}

	fmt.Println(url)

	bodyElement := element.NewElement(selectors.CSS, "#search_form_input_homepager")
	fmt.Printf("%v\n", bodyElement)

	// err = bodyElement.SendKeys(c, "Hello World")
	// if err != nil {
	// 	panic(err)
	// }

	// time.Sleep(5 * time.Second)

	// _, err = bodyElement.GetText(c)
	// if err != nil {
	// 	panic(err)
	// }

	// err = bodyElement.Clear(c)
	// if err != nil {
	// 	panic(err)
	// }

	// time.Sleep(5 * time.Second)

	// c.Refresh()

	// time.Sleep(5 * time.Second)

	// fmt.Println(text)

	err = bodyElement.WaitUntil(c, 60*time.Second).IsVisible()

	if err != nil {
		panic(err)
	}

}
