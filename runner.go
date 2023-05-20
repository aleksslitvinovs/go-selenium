package selenium

import (
	"fmt"
	"os"
	"sync"

	"github.com/aleksslitvinovs/go-selenium/logger"
	"github.com/fatih/color"
)

// TestFunction describes one test for the given session. It is used in
// selenium.Run() to execute tests.
type TestFunction func(s *Session)

type test struct {
	name     string
	fn       TestFunction
	s        *Session
	hadError bool
}

type runner struct {
	tests      []*test
	beforeAll  func()
	beforeEach TestFunction
	afterEach  TestFunction
	afterAll   func()
}

var r = &runner{}

// Run executes all set tests. If the client is not set, it sets one with the
// default driver based on the config settings.
func Run() {
	var errorCount int

	defer func() {
		MustStopClient()

		for _, t := range r.tests {
			if t.hadError {
				errorCount++
			}
		}

		if errorCount > 0 {
			os.Exit(1)
		}

		if errorCount > 0 {
			logger.Custom(color.RedString(
				"Failed! Success rate: %d/%d",
				len(r.tests)-errorCount, len(r.tests),
			))

			os.Exit(1)
		}

		logger.Custom(color.GreenString(
			"Passed! Success rate: %d/%d",
			len(r.tests)-errorCount, len(r.tests)),
		)
	}()

	if client == nil {
		err := SetClient(nil, nil)
		if err != nil {
			logger.Error(err)

			errorCount++

			return
		}
	}

	if config == nil {
		logger.Error("No configuration set")

		return
	}

	executeTests()
}

func executeTests() {
	runBeforeAll()

	pr := config.Runner.ParallelRuns
	if pr < 1 {
		pr = 1
	}

	fmt.Println("parallel", pr)

	jobs := make(chan *test, pr)
	defer close(jobs)

	wg := &sync.WaitGroup{}

	for i := 0; i < pr; i++ {
		go worker(jobs, wg)
	}

	for _, t := range r.tests {
		wg.Add(1)

		logger.Infof("running test: %s", t.name)

		jobs <- t
	}

	wg.Wait()

	runAfterAll()
}

func worker(tc <-chan *test, wg *sync.WaitGroup) {
	for t := range tc {
		runTest(t, wg)
	}
}

func runTest(t *test, wg *sync.WaitGroup) {
	defer wg.Done()

	defer handleTestPanic(t)

	s, err := NewSession()
	if err != nil {
		panic(err)
	}

	t.s = s

	client.ss.mu.Lock()
	client.ss.sessions[s] = true
	client.ss.mu.Unlock()

	defer s.DeleteSession()

	runBeforeEach(s)

	t.fn(t.s)

	runAfterEach(s)
}

func handleTestPanic(t *test) {
	err := recover()
	if err == nil {
		return
	}

	t.hadError = true

	if t.s == nil {
		return
	}

	switch v := err.(type) {
	case error:
		t.s.AddError(v.Error())
	case string:
		t.s.AddError(v)
	default:
		t.s.AddError(fmt.Sprintf("%v", v))
	}
}

func runBeforeAll() {
	if r.beforeAll != nil {
		r.beforeAll()
	}
}

func runBeforeEach(s *Session) {
	if r.beforeEach != nil {
		r.beforeEach(s)
	}
}

func runAfterEach(s *Session) {
	if r.afterEach != nil {
		r.afterEach(s)
	}
}

func runAfterAll() {
	if r.afterAll != nil {
		r.afterAll()
	}
}

// SetBeforeAll sets the function that will be executed before all tests.
func SetBeforeAll(fn func()) {
	r.beforeAll = fn
}

// SetBeforeEach sets the function that will be executed before each test.
func SetBeforeEach(fn TestFunction) {
	r.beforeEach = fn
}

// SetAfterEach sets the function that will be executed after each test.
func SetAfterEach(fn TestFunction) {
	r.afterEach = fn
}

// SetAfterAll sets the function that will be executed after all tests.
func SetAfterAll(fn func()) {
	r.afterAll = fn
}

// SetTest sets the test function. The name is used to identify the test is
// optional. If no name is provided, test_<test_id> is used. If the given name
// is already in use, test ID is appended to the name.
func SetTest(fn TestFunction, name ...string) {
	if len(name) == 0 {
		r.tests = append(r.tests, &test{
			name: fmt.Sprintf("test_%d", len(r.tests)),
			fn:   fn,
		})

		return
	}

	for _, t := range r.tests {
		if t.name == name[0] {
			r.tests = append(r.tests, &test{
				name: fmt.Sprintf("%s_%d", name[0], len(r.tests)),
				fn:   fn,
			})

			return
		}
	}

	r.tests = append(r.tests, &test{
		name: name[0],
		fn:   fn,
	})
}
