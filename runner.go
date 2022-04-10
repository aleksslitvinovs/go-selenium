package selenium

import (
	"fmt"
	"os"
	"sync"

	"github.com/theRealAlpaca/go-selenium/logger"
	"github.com/theRealAlpaca/go-selenium/types"
)

type runner struct {
	session     *Session
	tests       map[string]types.TestFunction
	testCounter int

	beforeAll  func()
	beforeEach func()
	afterEach  func()
	afterAll   func()
}

var Runner = &runner{tests: make(map[string]types.TestFunction)}

func Run() {
	if Runner.session == nil {
		panic("session is nil")
	}

	defer func() {
		MustStopClient()

		if len(Runner.session.errors) > 0 {
			os.Exit(1)
		}
	}()

	RunBeforeAll()

	wg := &sync.WaitGroup{}

	for n, t := range Runner.tests {
		wg.Add(1)
		RunBeforeEach()
		logger.Infof("running test: %s", n)

		go RunTest(t, wg)

		RunAfterEach()
	}

	wg.Wait()

	RunAfterAll()
}

func SetTest(fn types.TestFunction, name ...string) {
	if len(name) == 0 {
		Runner.tests[fmt.Sprintf("test_%d", Runner.testCounter)] = fn
		Runner.testCounter += 1

		return
	}

	Runner.tests[name[0]] = fn
}

func SetTests(fns map[string]types.TestFunction) {
	for n, fn := range fns {
		SetTest(fn, n)
	}
}

func RunBeforeAll() {
	if Runner.beforeAll != nil {
		Runner.beforeAll()
	}
}

func RunBeforeEach() {
	if Runner.beforeEach != nil {
		Runner.beforeEach()
	}
}

func RunTest(fn types.TestFunction, wg *sync.WaitGroup) {
	defer wg.Done()

	defer func() {
		err := recover()
		if err == nil {
			return
		}

		v, ok := err.(string)
		if !ok {
			Runner.session.AddError("unknown error occurred")

			return
		}

		Runner.session.AddError(v)
	}()

	fn(Runner.session)
}

func RunAfterEach() {
	if Runner.afterEach != nil {
		Runner.afterEach()
	}
}

func RunAfterAll() {
	if Runner.afterAll != nil {
		Runner.afterAll()
	}
}
