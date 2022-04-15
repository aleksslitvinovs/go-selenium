package selenium

import (
	"fmt"
	"os"
	"sync"

	"github.com/theRealAlpaca/go-selenium/config"
	"github.com/theRealAlpaca/go-selenium/logger"
	"github.com/theRealAlpaca/go-selenium/types"
)

type runner struct {
	tests       map[string]types.TestFunction
	testCounter int
	hadErrors   bool

	beforeAll  func()
	beforeEach func()
	afterEach  func()
	afterAll   func()
}

var r = &runner{tests: make(map[string]types.TestFunction)}

func Run() {
	defer func() {
		MustStopClient()

		if r.hadErrors {
			os.Exit(1)
		}
	}()

	if Client == nil || config.Config == nil {
		logger.Error("Client is not set")

		r.hadErrors = true

		return
	}

	runBeforeAll()

	pr := config.Config.Runner.ParallelRuns
	if pr < 1 {
		pr = 1
	}

	jobs := make(chan types.TestFunction, pr)
	defer close(jobs)

	wg := &sync.WaitGroup{}

	for i := 0; i < pr; i++ {
		go worker(jobs, wg)
	}

	for n, t := range r.tests {
		wg.Add(1)

		runBeforeEach()

		logger.Infof("running test: %s", n)
		jobs <- t

		runAfterEach()
	}

	wg.Wait()

	runAfterAll()
}

func worker(tf <-chan types.TestFunction, wg *sync.WaitGroup) {
	for t := range tf {
		runTest(t, wg)
	}
}

func runTest(fn types.TestFunction, wg *sync.WaitGroup) {
	var s types.Sessioner

	defer wg.Done()

	defer func() {
		s.DeleteSession()

		err := recover()
		if err == nil {
			return
		}

		r.hadErrors = true

		switch v := err.(type) {
		case error:
			s.AddError(v.Error())
		case string:
			s.AddError(v)
		default:
			s.AddError(fmt.Sprintf("%v", v))
		}
	}()

	s, err := CreateSession()
	if err != nil {
		panic(err)
	}

	Client.sessions[s.(*Session)] = true

	fn(s)
}

func runBeforeAll() {
	if r.beforeAll != nil {
		r.beforeAll()
	}
}

func runBeforeEach() {
	if r.beforeEach != nil {
		r.beforeEach()
	}
}

func runAfterEach() {
	if r.afterEach != nil {
		r.afterEach()
	}
}

func runAfterAll() {
	if r.afterAll != nil {
		r.afterAll()
	}
}

func SetBeforeAll(f func()) {
	r.beforeAll = f
}

func SetBeforeEach(f func()) {
	r.beforeEach = f
}

func SetAfterEach(f func()) {
	r.afterEach = f
}

func SetAfterAll(f func()) {
	r.afterAll = f
}

func SetTest(fn types.TestFunction, name ...string) {
	defer func() {
		r.testCounter++
	}()

	if len(name) == 0 {
		r.tests[fmt.Sprintf("test_%d", r.testCounter)] = fn

		return
	}

	if r.tests[name[0]] != nil {
		r.tests[fmt.Sprintf("%s_%d", name[0], r.testCounter)] = fn

		return
	}

	r.tests[name[0]] = fn
}

func SetTests(fns map[string]types.TestFunction) {
	for n, fn := range fns {
		SetTest(fn, n)
	}
}
