package selenium_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/theRealAlpaca/go-selenium"
	"github.com/theRealAlpaca/go-selenium/selector"
)

func Test(t *testing.T) {
	// d, err := selenium.NewDriver(
	// 	"/Users/aleksslitvinovs/Downloads/chromedriver",
	// 	"http://localhost:4444",
	// )
	// if err != nil {
	// 	panic(err)
	// }
	err := selenium.SetClient(nil, nil)
	if err != nil {
		panic(err)
	}

	selenium.SetTest(AssertTest)
	// selenium.SetTest(IFrameTest)
	// selenium.SetTest(JitsiTest)
	// selenium.SetTest(MyTest)

	selenium.Run()
}

func AssertTest(s *selenium.Session) {
	s.OpenURL("https://duckduckgo.com/")
	s.TakeScreeshot("test.jpeg")

	fmt.Println(
		s.NewElement(
			&selenium.E{"#search_form_input_homepage", selector.CSS},
		).GetAttribute("id"),
	)

	s.NewElement("#search_form_input_homepage").
		ShouldHave().Attribute("id").
		EqualsTo("search_form_input_homepage")
	// s.NewElement(".text_promo--text").ShouldHave().Text().EndsWith("Beta")
	fmt.Println(s.NewElement("text_promo--text").GetText())
	fmt.Println(
		"testing value",
		s.NewElement("text_promo--text").GetAttribute("value"),
	)
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

func MyTest(s *selenium.Session) {
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
