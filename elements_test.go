package selenium_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aleksslitvinovs/go-selenium"
	"github.com/aleksslitvinovs/go-selenium/keys"
)

func TestElements(t *testing.T) {
	selenium.SetTest(MultipleElementsTest)

	selenium.Run()
}

func MultipleElementsTest(s *selenium.Session) {
	s.OpenURL("https://duckduckgo.com/")

	s.NewElement("input[type=text]").
		WaitFor(10 * time.Second).UntilIsVisible().
		SendKeys("WebDriver").
		SendKeys(keys.Enter)

	titles := s.NewElements("[data-testid=result-title-a]")

	fmt.Println("elements size", titles.Size())
	fmt.Println("elements ids", titles.Elements())

	for _, e := range titles.Elements() {
		fmt.Println("element text", e.GetText())
	}
}
