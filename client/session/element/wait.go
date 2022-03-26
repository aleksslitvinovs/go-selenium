package element

import (
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/client/session"
	"github.com/theRealAlpaca/go-selenium/util"
)

type waiter struct {
	e       *Element
	timeout time.Duration
}

func (e *Element) WaitFor(timeout time.Duration) *waiter {
	return &waiter{
		e:       e,
		timeout: timeout,
	}
}

func (w *waiter) UntilIsVisible() *Element {
	return waitCondition(w, w.e.IsVisible, true, "isible")
}
func (w *waiter) UntilIsNotVisible() *Element {
	return waitCondition(w, w.e.IsVisible, false, "not visible")
}

func (w *waiter) UntilIsEnabled() *Element {
	return waitCondition(w, w.e.IsEnabled, true, "enabled")
}

func (w *waiter) UntilIsNotEnabled() *Element {
	return waitCondition(w, w.e.IsEnabled, false, "not enabled")
}

func (w *waiter) UntilIsSelected() *Element {
	return waitCondition(w, w.e.IsSelected, true, "selected")
}

func (w *waiter) UntilIsNotSelected() *Element {
	return waitCondition(w, w.e.IsSelected, false, "not selected")
}

func waitCondition(
	w *waiter,
	condition func(s *session.Session) (bool, error),
	expected bool,
	conditionName string,
) *Element {
	startTime := time.Now()
	endTime := startTime.Add(w.timeout)

	w.e.IgnoreNotFound = true

	defer func() {
		w.e.IgnoreNotFound = false
	}()

	for {
		if endTime.Before(time.Now()) {
			util.HandleError(
				w.e.Session, errors.Errorf(
					"Element %q is not %s after %s (time elapsed %s)",
					w.e.Selector,
					conditionName,
					w.timeout,
					time.Since(startTime),
				),
			)

			return w.e
		}

		actual, err := condition(w.e.Session)
		if err != nil {
			util.HandleError(
				w.e.Session, errors.Wrap(err, "could not get condition"),
			)
		}

		if actual == expected {
			return w.e
		}

		time.Sleep(time.Second)
	}
}
