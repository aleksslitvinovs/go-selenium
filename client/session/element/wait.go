package element

import (
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/client/session"
)

type waiter struct {
	s       *session.Session
	e       *Element
	timeout time.Duration
}

func (e *Element) WaitFor(s *session.Session, timeout time.Duration) *waiter {
	return &waiter{
		s:       s,
		e:       e,
		timeout: timeout,
	}
}

func (w *waiter) UntilIsVisible() error {
	return waitCondition(w, w.e.IsVisible, true)
}
func (w *waiter) UntilIsNotVisible() error {
	return waitCondition(w, w.e.IsVisible, false)
}

func (w *waiter) UntilIsEnabled() error {
	return waitCondition(w, w.e.IsEnabled, true)
}

func (w *waiter) UntilIsNotEnabled() error {
	return waitCondition(w, w.e.IsEnabled, false)
}

func (w *waiter) UntilIsSelected() error {
	return waitCondition(w, w.e.IsSelected, true)
}

func (w *waiter) UntilIsNotSelected() error {
	return waitCondition(w, w.e.IsSelected, false)
}

func waitCondition(
	w *waiter, condition func(s *session.Session) (bool, error), expected bool,
) error {
	startTime := time.Now()
	endTime := startTime.Add(w.timeout)

	for {
		if endTime.Before(time.Now()) {
			return errors.New("exceeded timeout to reach expected condition")
		}

		actual, err := condition(w.s)
		if err != nil {
			return errors.Wrap(err, "could not get condition")
		}

		if actual == expected {
			return nil
		}
	}
}
