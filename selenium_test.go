package selenium_test

import (
	"testing"
	"time"

	"github.com/theRealAlpaca/go-selenium"
)

func Test(t *testing.T) {
	d, err := selenium.NewDriver(
		"/Users/aleksslitvinovs/Downloads/chromedriver",
		"http://localhost:4444",
	)
	if err != nil {
		panic(err)
	}

	c, err := selenium.NewClient(d, nil)
	if err != nil {
		panic(err)
	}

	defer c.MustStop()

	s, err := c.CreateSession()
	if err != nil {
		panic(err)
	}

	// s.OpenURL("https://app.stage.loadero.com/login")
	s.OpenURL("https://app.stage.loadero.com/login")
	s.NewElement("#email")
	// s.NewWindow()

	// fmt.Println(s.GetWindowHandle())
	// fmt.Println(s.GetWindowHandles())
	// s.SwitchToParentFrame()
	// // s.SwitchFrame(1)

	// fmt.Println(s.GetWindowHandle())

	// s.CloseWindow()
	// s.
	// 	NewElement(".account__logo-link").
	// 	WaitFor(time.Second * 5).
	// 	UntilIsVisible().
	// 	Click()

	// s.Back()
	// fmt.Println(s.GetTitle())

	// s.Forward()
	// fmt.Println(s.GetTitle())

	s.NewElement(".sign-in-form")
	s.NewElement("#username").
		WaitFor(time.Second * 5).UntilIsVisible().
		SendKeys("testing@loadero.abc")

	s.NewElement("#password").
		WaitFor(time.Second * 5).UntilIsVisible().
		SendKeys("password")

	s.NewElement(".button--primary").
		WaitFor(time.Second * 5).UntilIsVisible().
		Click()

	time.Sleep(10 * time.Second)
}
