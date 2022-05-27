package selenium_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/theRealAlpaca/go-selenium"
)

const url = "https://meet.jit.si/RTU_test"

func Test(t *testing.T) {
	selenium.SetAfterEach(func(s *selenium.Session) {
		fmt.Println(time.Now().Second())
	})
	selenium.SetTest(MyTest)
	// selenium.SetTest(JitsiTest1)
	// selenium.SetTest(JitsiTest2)

	selenium.Run()
}

func JitsiTest1(s *selenium.Session) {
	time.Sleep(time.Second * 5)

	fmt.Println("First done")
	s.OpenURL(url)

	s.NewElement(`[aria-label="Join meeting"]`).
		WaitFor(time.Second * 5).UntilIsVisible().
		Click()

	s.NewElement(".prejoin-input-area input").
		WaitFor(time.Second * 5).UntilIsVisible().
		SendKeys("Jitsi 1")

	time.Sleep(10 * time.Second)

	s.TakeScreenshot("jitsi.png")
}

func JitsiTest2(s *selenium.Session) {
	time.Sleep(time.Second * 6)

	fmt.Println("Second done")
	s.OpenURL(url)

	s.NewElement(`[aria-label="Join meeting"]`).
		WaitFor(time.Second * 5).UntilIsVisible().
		Click()

	s.NewElement(".prejoin-input-area input").
		WaitFor(time.Second * 5).UntilIsVisible().
		SendKeys("Jitsi 2")

	time.Sleep(10 * time.Second)

	s.TakeScreenshot("jitsi.png")
}
