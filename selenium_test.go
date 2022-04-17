package selenium_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/theRealAlpaca/go-selenium"
	"github.com/theRealAlpaca/go-selenium/key"
)

func Test(t *testing.T) {
	selenium.SetTest(AssertTest)

	selenium.Run()
}

func MyTest(s *selenium.Session) {
	s.OpenURL("https://duckduckgo.com/")

	s.NewElement("#search_form_input_homepage").
		WaitFor(10 * time.Second).UntilIsVisible().
		SendKeys("WebDriver").
		SendKeys(key.Enter)

	result := s.NewElement("#r1-0 .result__a").
		WaitFor(10 * time.Second).UntilIsVisible().
		GetText()

	fmt.Printf("DuckDuckGo result: %s\n", result)
}

func Test2(t *testing.T) {
	// d, err := selenium.NewDriver(
	// 	"/Users/aleksslitvinovs/Downloads/chromedriver",
	// 	"http://localhost:4444",
	// )
	// if err != nil {
	// 	panic(err)
	// }

	// selenium.SetTest(AssertTest)
	// selenium.SetTest(IFrameTest)
	// selenium.SetTest(JitsiTest)
	// selenium.SetTest(Testy)
}

func AssertTest(s *selenium.Session) {
	s.OpenURL("https://duckduckgo.com/")

	s.NewElement("#search_form_input_homepage").
		WaitFor(10 * time.Second).UntilIsVisible().
		SendKeys("WebDriver").
		SendKeys(key.Enter)

	e := s.NewElement("#search_form_input").
		WaitFor(10 * time.Second).UntilIsVisible()

	e.ShouldHave().Attribute("value").EqualTo("WebDriver")
	e.ShouldHave().Attribute("value").Not().EqualTo("WebDriver_fail")

	e.ShouldHave().Attribute("value").EqualTo("WebDriver_fail")
	e.ShouldHave().Attribute("value").Not().EqualTo("WebDriver")

	result := s.NewElement("#r1-0 .result__a").
		WaitFor(10 * time.Second).UntilIsVisible().
		GetText()

	fmt.Printf("DuckDuckGo result: %s\n", result)
}

func IFrameTest(s *selenium.Session) {
	s.OpenURL("https://jsfiddle.net/westonruter/6mSuK/")

	t := s.NewElement(".iframeCont")
	t.WaitFor(10 * time.Second).UntilIsVisible()

	s.SwitchToFrame(s.NewElement(`[name="result"]`)).
		SwitchToFrame((s.NewElement("iframe")))

	fmt.Println(
		"Result",
		s.NewElement("#ca-nstab-project").
			WaitFor(10*time.Second).UntilIsVisible().
			GetText(),
	)
	s.SwitchToParentFrame()

	fmt.Println(
		"Result 2",
		s.NewElement("body").
			WaitFor(10*time.Second).UntilIsVisible().
			GetText(),
	)

	s.SwitchToFrame(nil)
	fmt.Println(
		"Result 3",
		s.NewElement(".profileDetails .company").
			WaitFor(10*time.Second).UntilIsVisible().
			GetText(),
	)
}

func JitsiTest(s *selenium.Session) {
	s.OpenURL("https://meet.jit.si/LoaderoWebRTC_R")

	s.NewElement(`[aria-label="Join meeting"]`).
		WaitFor(time.Second * 5).UntilIsVisible().
		Click()

	time.Sleep(10 * time.Second)
}

func Testy(s *selenium.Session) {
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
