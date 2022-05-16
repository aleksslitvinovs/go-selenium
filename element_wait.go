package selenium

import (
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/logger"
)

// Waiter is a helper struct to wait for an element to be present, visible, etc.
type Waiter struct {
	e       *Element
	timeout time.Duration
}

// WaitFor creates an instance of *Waiter with the provided timeout duration.
func (e *Element) WaitFor(timeout time.Duration) *Waiter {
	return &Waiter{
		e:       e,
		timeout: timeout,
	}
}

// UntilIsPresent waits until the element is present.
func (w *Waiter) UntilIsPresent() *Element {
	return waitPresent(w, true)
}

// UntilIsNotPresent waits until the element is not present.
func (w *Waiter) UntilIsNotPresent() *Element {
	return waitPresent(w, false)
}

// UntilIsVisible waits until the element is visible.
func (w *Waiter) UntilIsVisible() *Element {
	return waitCondition(w, w.e.isVisible, true, "visible")
}

// UntilIsNotVisible waits until the element is visible.
func (w *Waiter) UntilIsNotVisible() *Element {
	return waitCondition(w, w.e.isVisible, false, "not visible")
}

// UntilIsEnabled waits until the element is enabled.
func (w *Waiter) UntilIsEnabled() *Element {
	return waitCondition(w, w.e.isEnabled, true, "enabled")
}

// UntilIsNotEnabled waits until the element is not enabled.
func (w *Waiter) UntilIsNotEnabled() *Element {
	return waitCondition(w, w.e.isEnabled, false, "not enabled")
}

// UntilIsSelected waits until the element is selected.
func (w *Waiter) UntilIsSelected() *Element {
	return waitCondition(w, w.e.isSelected, true, "selected")
}

// UntilIsNotSelected waits until the element is not selected.
func (w *Waiter) UntilIsNotSelected() *Element {
	return waitCondition(w, w.e.isSelected, false, "not selected")
}

func waitCondition(
	w *Waiter,
	condition func() (*response, error),
	expected bool,
	conditionName string,
) *Element {
	startTime := time.Now()
	endTime := startTime.Add(w.timeout)

	intialSettings := *w.e.settings

	w.e.settings.IgnoreNotFound = true

	defer func() {
		w.e.settings = &intialSettings
	}()

	for {
		if endTime.Before(time.Now()) {
			handleError(
				nil,
				errors.Errorf(
					"Element %q is not %s after %s (time elapsed %dms)",
					w.e.Selector,
					conditionName,
					w.timeout,
					time.Since(startTime).Milliseconds(),
				),
			)

			return w.e
		}

		res, err := condition()
		if err != nil {
			handleError(res, err)

			return w.e
		}

		if res.Value != nil {
			if v, ok := res.Value.(bool); ok {
				if v == expected {
					logger.Infof(
						"Element %q is %s after %s (time elapsed %dms)",
						w.e.Selector,
						conditionName,
						w.timeout,
						time.Since(startTime).Milliseconds(),
					)

					return w.e
				}
			}
		}

		time.Sleep(w.e.settings.PollInterval.Duration)
	}
}

func waitPresent(
	w *Waiter,
	bePresent bool,
) *Element {
	startTime := time.Now()
	endTime := startTime.Add(w.timeout)

	intialSettings := *w.e.settings

	w.e.settings.IgnoreNotFound = true

	defer func() {
		w.e.settings = &intialSettings
	}()

	conditionName := "present"
	if !bePresent {
		conditionName = "not present"
	}

	for {
		if endTime.Before(time.Now()) {
			handleError(
				nil,
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

		id, err := w.e.findElement()
		if err != nil {
			handleError(nil, err)

			return w.e
		}

		if id != "" && bePresent || id == "" && !bePresent {
			logger.Infof(
				"Element %q is %s after %s (time elapsed %dms)",
				w.e.Selector,
				conditionName,
				w.timeout,
				time.Since(startTime).Milliseconds(),
			)

			w.e.id = id

			return w.e
		}

		time.Sleep(w.e.settings.PollInterval.Duration)
	}
}
