package selenium_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium"
	"github.com/theRealAlpaca/go-selenium/key"
)

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

func EbayTestSearch(s *selenium.Session) {
	s.OpenURL("https://www.amazon.com/")

	s.NewElement(".nav-search-field .nav-input").
		WaitFor(10 * time.Second).UntilIsVisible().
		SendKeys("MacBook").
		SendKeys(key.Enter)

	fmt.Println(
		"First result text:",
		s.NewElement(`[cel_widget_id="MAIN-SEARCH_RESULTS-2]`).
			WaitFor(10*time.Second).UntilIsVisible().
			GetText(),
	)

	s.TakeScreenshot("first result.png")
}

func UltimateQAForm(s *selenium.Session) {
	s.OpenURL("https://ultimateqa.com/filling-out-forms/")

	s.NewElement("#et_pb_contact_name_1").
		WaitFor(10 * time.Second).UntilIsVisible().
		SendKeys("John Smith")

	s.NewElement("#et_pb_contact_message_1").
		WaitFor(10 * time.Second).UntilIsVisible().
		SendKeys("Hello, I'm John Smith!")

	fillCatpcha(s)

	s.NewElement("#et_pb_contact_form_1 .et_pb_button").
		WaitFor(10 * time.Second).UntilIsVisible().
		Click()

	fillCatpcha(s)

	s.NewElement("#et_pb_contact_form_1 .et_pb_button").
		WaitFor(10 * time.Second).UntilIsVisible().
		Click()

	time.Sleep(5 * time.Second)

	s.NewElement("#et_pb_contact_form_1 .et-pb-contact-message p").
		WaitFor(10 * time.Second).UntilIsVisible().
		ShouldHave().Text().EqualTo("Thanks for contacting us")
}

func fillCatpcha(s *selenium.Session) {
	captcha := s.NewElement(".et_pb_contact_captcha").
		WaitFor(10 * time.Second).UntilIsVisible()

	firstDigit, err := strconv.Atoi(captcha.GetAttribute("data-first_digit"))
	if err != nil {
		panic(errors.Wrap(err, "could not get first digit"))
	}

	secondDigit, err := strconv.Atoi(captcha.GetAttribute("data-second_digit"))
	if err != nil {
		panic(errors.Wrap(err, "could not get second digit"))
	}

	captcha.SendKeys(strconv.Itoa(firstDigit + secondDigit))

	fmt.Println(firstDigit + secondDigit)
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

func DuckDuckGo(s *selenium.Session) {
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
