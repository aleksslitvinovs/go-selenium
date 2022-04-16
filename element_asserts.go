package selenium

import (
	"strings"

	"github.com/pkg/errors"
)

type Asserter struct {
	e *Element
}

type Valuer struct {
	actual string
	e      *Element
}

func (e *Element) ShouldHave() *Asserter {
	return &Asserter{
		e: e,
	}
}

func (a *Asserter) Text() *Valuer {
	text := a.e.GetText()

	return &Valuer{
		actual: text,
		e:      a.e,
	}
}

func (a *Asserter) Attribute(attribute string) *Valuer {
	return &Valuer{
		actual: a.e.GetAttribute(attribute),
		e:      a.e,
	}
}

func (t *Valuer) EqualsTo(expected string) {
	if t.actual != expected {
		handleError(
			nil,
			errors.Errorf(
				"text %q should have been equal to %q", t.actual, expected,
			),
		)
	}
}

func (t *Valuer) NotEqualsTo(expected string) {
	if t.actual == expected {
		handleError(
			nil,
			errors.Errorf(
				"text %q should have not been equal to %q", t.actual, expected,
			),
		)
	}
}

func (t *Valuer) StartsWith(expected string) {
	if strings.HasPrefix(t.actual, expected) {
		handleError(
			nil,
			errors.Errorf(
				"text %q should have started with %q", t.actual, expected,
			),
		)
	}
}

func (t *Valuer) EndsWith(expected string) {
	if !strings.HasSuffix(t.actual, expected) {
		handleError(
			nil,
			errors.Errorf(
				"text %q should have ended with %q", t.actual, expected,
			),
		)
	}
}

func (t *Valuer) Contains(expected string) {
	if !strings.Contains(t.actual, expected) {
		handleError(
			nil,
			errors.Errorf(
				"text %q shoud have contained %q", t.actual, expected,
			),
		)
	}
}

func (t *Valuer) NotContains(expected string) {
	if strings.Contains(t.actual, expected) {
		handleError(
			nil,
			errors.Errorf(
				"text %q should not have contained %q", t.actual, expected,
			),
		)
	}
}
