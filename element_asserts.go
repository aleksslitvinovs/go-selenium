package selenium

import (
	"fmt"
	"strings"

	"github.com/aleksslitvinovs/go-selenium/logger"
	"github.com/pkg/errors"
)

// Asserter is a helper struct to assert the element's text, attributes, etc.
type Asserter struct {
	e *Element
}

// Valuer is a helper struct to compare the actual value of an element with the
// expected using various comparison functions.
type Valuer struct {
	actual   string
	e        *Element
	property string
	isEqual  bool
}

// ShouldHave returns an Asserter for the given element. The returned Asserter
// can be used to assert the element's text, attributes, etc.
func (e *Element) ShouldHave() *Asserter {
	return &Asserter{
		e: e,
	}
}

// Text allows asserting the element's text.
func (a *Asserter) Text() *Valuer {
	text := a.e.GetText()

	return &Valuer{
		actual:   text,
		e:        a.e,
		property: "text",
		isEqual:  true,
	}
}

// Attribute allows asserting the element's attributes.
func (a *Asserter) Attribute(attribute string) *Valuer {
	return &Valuer{
		actual:   a.e.GetAttribute(attribute),
		e:        a.e,
		property: fmt.Sprintf("attribute %q", attribute),
		isEqual:  true,
	}
}

// Not negates the following assertion.
func (v *Valuer) Not() *Valuer {
	v.isEqual = false

	return v
}

type comparer struct {
	past    string
	present string
}

var (
	equalTo   = comparer{"been equal to", "equal"}
	startWith = comparer{"started with", "start with"}
	endWith   = comparer{"ended with", "end with"}
	contain   = comparer{"contained", "contain"}
)

// EqualTo asserts that the actual value of the element is equal to the given.
func (v *Valuer) EqualTo(expected string) {
	v.compare(expected, equalTo)
}

// StartWith asserts that the actual value of the element starts with the
// given.
func (v *Valuer) StartWith(expected string) {
	v.compare(expected, startWith)
}

// EndWith asserts that the actual value of the element ends with the given.
func (v *Valuer) EndWith(expected string) {
	v.compare(expected, endWith)
}

// Contain asserts that the actual value of the element contains the given.
func (v *Valuer) Contain(expected string) {
	v.compare(expected, contain)
}

func (v *Valuer) compare(expected string, cmp comparer) {
	var result bool

	switch cmp {
	case equalTo:
		result = v.actual == expected
	case startWith:
		result = strings.HasPrefix(v.actual, expected)
	case endWith:
		result = strings.HasSuffix(v.actual, expected)
	case contain:
		result = strings.Contains(v.actual, expected)
	}

	not := comparer{"", "does"}
	if !v.isEqual {
		not = comparer{"not ", "does not"}
	}

	if xnor(result, v.isEqual) {
		logger.Infof(
			"element's %s %s %s %q",
			v.property, not.present, cmp.present, expected,
		)

		return
	}

	handleError(
		nil,
		errors.Errorf(
			"element's %s should have %s%s %q, actual value %q",
			v.property, not.past, cmp.past, expected, v.actual,
		),
	)
}

func xnor(a, b bool) bool {
	return (a && b) || (!a && !b)
}
