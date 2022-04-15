package selenium_test

import (
	"testing"
	"time"

	"github.com/theRealAlpaca/go-selenium"
	"github.com/theRealAlpaca/go-selenium/types"
)

func Test(t *testing.T) {
	// d, err := selenium.NewDriver(
	// 	"/Users/aleksslitvinovs/Downloads/chromedriver",
	// 	"http://localhost:4444",
	// )
	// if err != nil {
	// 	panic(err)
	// }

	_, err := selenium.StartClient(nil, nil)
	if err != nil {
		panic(err)
	}

	selenium.SetTest(JitsiTest)
	selenium.SetTest(MyTest)

	selenium.Run()
}

func JitsiTest(s types.Sessioner) {
	s.OpenURL("https://meet.jit.si/LoaderoWebRTC_R")

	s.NewElement(`[aria-label="Join meeting"]`).
		WaitFor(time.Second * 5).UntilIsVisible().
		Click()

	time.Sleep(10 * time.Second)
}

func MyTest(s types.Sessioner) {
	s.OpenURL("https://app.stage.loadero.com/login")
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
}
