# go-selenium

[![Go Reference](https://pkg.go.dev/badge/github.com/theRealAlpaca/go-selenium.svg)](https://pkg.go.dev/github.com/theRealAlpaca/go-selenium)
[![Go Report Card](https://goreportcard.com/badge/github.com/theRealAlpaca/go-selenium)](https://goreportcard.com/report/github.com/theRealAlpaca/go-selenium)

## About

go-selenium is a library for interacting with the
[browserdriver](https://www.selenium.dev/documentation/overview/components/#terminology)
by implementing [W3C WebDriver recomendation](https://www.w3.org/TR/webdriver1/)
using the Go programming language.

go-selenium supports the following browsers:

- Chrome ([ChromeDriver](https://chromedriver.chromium.org/home))
- Firefox([GeckoDriver](https://github.com/mozilla/geckodriver))

## Installation

```bash
go get github.com/theRealAlpaca/go-selenium
```

## Usage

Create `selenium_test.go` file with the following example code:

```go
package selenium_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/theRealAlpaca/go-selenium"
	"github.com/theRealAlpaca/go-selenium/key"
)

func Test(t *testing.T) {
	selenium.SetTest(MyTest)

	selenium.Run()
}

func MyTest(s *selenium.Session) {
	s.OpenURL("https://duckduckgo.com/")

	s.NewElement("#search_form_input_homepage").
		WaitFor(10 * time.Second).UntilIsVisible().
		SendKeys("theRealAlpaca/go-selenium").
		SendKeys(key.Enter)

	result := s.NewElement("#r1-0 [data-testid=result-title-a]").
		WaitFor(10 * time.Second).UntilIsVisible().
		GetText()

	fmt.Printf("DuckDuckGo result: %s\n", result)
}
```

To run the, run `go test`.

## Configuration

Even though the library is designed to work with the default configuration, it
has many configurable options that are defined in `.goseleniumrc.json` file
(generated automatically on the first run).

| Option                       | Description                                                                | Type                              | Default                   |
| ---------------------------- | -------------------------------------------------------------------------- | --------------------------------- | ------------------------- |
| `logging`                    | Logging level                                                              | `string`                          | `"info"`                  |
| `soft_asserts`               | Use soft assertions, i.e., continue executing the test in case of an error | `bool`                            | `true`                    |
| `screenshot_dir`             | Directory in which save screenshots                                        | `string`                          | `""`                      |
| `raise_errors_automatically` | Raise errors automatically when the test ends                              | `bool`                            | `true`                    |
| `runner`                     |                                                                            | `object`                          |                           |
| `runner.parallel_runs`       | Number of parallel tests to execute                                        | `int`                             | `1`                       |
| `element`                    |                                                                            | `object`                          |                           |
| `element.selector_type`      | Default selector type used when locating element                           | `string`                          | `css selector`            |
| `element.ignore_not_found`   | Throw error if element is not found                                        | `bool`                            | `false`                   |
| `element.retry_timeout`      | Timeout for trying to locate the given element                             | [`time`](./README.md#time-format) | `10s`                     |
| `element.poll_interval`      | Time interval to validate element's state when using `WaitFor()` command   | [`time`](./README.md#time-format) | `500ms`                   |
| `webdriver`                  |                                                                            | `object`                          |                           |
| `webdriver.browser`          | Browser to use                                                             | `string`                          | `"chrome"`                |
| `webdriver.manual_start`     | Start browser driver process manually                                      | `bool`                            | `false`                   |
| `webdriver.binary_path`      | Path to browser driver binary                                              | `string`                          | `"./chromedriver"`        |
| `webdriver.remote_url`       | URL to which WebDriver commands are sent                                   | `string`                          | `"http://localhost:4444"` |
| `webdriver.timeout`          | Time which which browser driver should be ready to accept command          | [`time`](./README.md#time-format) | `"10 s"`                  |
| `webdriver.capabilities`     | Browser capabilities                                                       | `map[string]interface{}`          | `{}`                      |

### Time format

Time type is defined as a string, e.g., `500ms`, with the following format:

```
<duration><unit>
```

where `<duration>` is an integer and `<unit>` is one of the following time units:

- `ns` - nanoseconds
- `us` (or `µs`) - microseconds
- `ms` - milliseconds
- `s` - seconds
- `m` - minutes
- `h` - hours
