package main

import (
	"time"

	"github.com/theRealAlpaca/go-selenium"
)

//nolint:errcheck
func main() {
	c, s := selenium.Start(nil)

	defer c.Stop()

	// s.OpenURL("https://app.stage.loadero.com/login")
	s.Navigation.OpenURL("https://app.stage.loadero.com/login")
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

	s.NewElement(".sign-in-form").WaitFor(time.Second * 5).UntilIsVisible()
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
