package selenium

import (
	"strings"

	"github.com/pkg/errors"
)

type asserter struct {
	e *element
}

type valuer struct {
	actual string
	e      *element
}

func (e *element) ShouldHave() *asserter {
	return &asserter{
		e: e,
	}
}

func (a *asserter) Text() *valuer {
	text := a.e.GetText()

	return &valuer{
		actual: text,
		e:      a.e,
	}
}

func (a *asserter) Attribute(attribute string) *valuer {
	return &valuer{
		actual: a.e.GetAttribute(attribute),
		e:      a.e,
	}
}

func (t *valuer) EqualsTo(expected string) {
	if t.actual != expected {
		handleError(
			nil,
			errors.Errorf(
				"text %q should have been equal to %q", t.actual, expected,
			),
		)
	}
}

func (t *valuer) NotEqualsTo(expected string) {
	if t.actual == expected {
		handleError(
			nil,
			errors.Errorf(
				"text %q should have not been equal to %q", t.actual, expected,
			),
		)
	}
}

func (t *valuer) StartsWith(expected string) {
	if strings.HasPrefix(t.actual, expected) {
		handleError(
			nil,
			errors.Errorf(
				"text %q should have started with %q", t.actual, expected,
			),
		)
	}
}

func (t *valuer) EndsWith(expected string) {
	if !strings.HasSuffix(t.actual, expected) {
		handleError(
			nil,
			errors.Errorf(
				"text %q should have ended with %q", t.actual, expected,
			),
		)
	}
}

func (t *valuer) Contains(expected string) {
	if !strings.Contains(t.actual, expected) {
		handleError(
			nil,
			errors.Errorf(
				"text %q shoud have contained %q", t.actual, expected,
			),
		)
	}
}

func (t *valuer) NotContains(expected string) {
	if strings.Contains(t.actual, expected) {
		handleError(
			nil,
			errors.Errorf(
				"text %q should not have contained %q", t.actual, expected,
			),
		)
	}
}
