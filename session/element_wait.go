package session

import (
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
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
	return waitCondition(w, w.e.isVisible, true, "is visible")
}
func (w *waiter) UntilIsNotVisible() *Element {
	return waitCondition(w, w.e.isVisible, false, "not visible")
}

func (w *waiter) UntilIsEnabled() *Element {
	return waitCondition(w, w.e.isEnabled, true, "enabled")
}

func (w *waiter) UntilIsNotEnabled() *Element {
	return waitCondition(w, w.e.isEnabled, false, "not enabled")
}

func (w *waiter) UntilIsSelected() *Element {
	return waitCondition(w, w.e.isSelected, true, "selected")
}

func (w *waiter) UntilIsNotSelected() *Element {
	return waitCondition(w, w.e.isSelected, false, "not selected")
}

func waitCondition(
	w *waiter,
	condition func() (*api.Response, error),
	expected bool,
	conditionName string,
) *Element {
	startTime := time.Now()
	endTime := startTime.Add(w.timeout)

	intialSettings := *w.e.Settings

	w.e.Settings.IgnoreNotFound = true

	defer func() {
		w.e.Settings = &intialSettings
	}()

	for {
		if endTime.Before(time.Now()) {
			util.HandleError(
				w.e.Session,
				errors.Errorf(
					"Element %q is not %s after %s (time elapsed %s)",
					w.e.Selector,
					conditionName,
					w.timeout,
					time.Since(startTime),
				),
			)

			return w.e
		}

		res, err := condition()
		if err != nil {
			if errors.As(err, &ErrWebIDNotSet) {
				time.Sleep(w.e.Settings.PollInterval.Duration)

				continue
			}

			util.HandleError(
				w.e.Session, errors.Wrap(err, "could not get condition"),
			)

			return w.e
		}

		if res.Value.(bool) == expected {
			return w.e
		}

		time.Sleep(w.e.Settings.PollInterval.Duration)
	}
}
