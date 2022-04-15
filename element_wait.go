package selenium

import (
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/logger"
)

type waiter struct {
	we      *Element
	timeout time.Duration
}

var (
	ErrWebIDNotSet = errors.New("WebID not set")
)

func (we *Element) WaitFor(timeout time.Duration) *waiter {
	return &waiter{
		we:      we,
		timeout: timeout,
	}
}

func (w *waiter) UntilIsPresent() *Element {
	return waitPresent(w, true)
}
func (w *waiter) UntilIsNotPresent() *Element {
	return waitPresent(w, false)
}

func (w *waiter) UntilIsVisible() *Element {
	return waitCondition(w, w.we.isVisible, true, "visible")
}
func (w *waiter) UntilIsNotVisible() *Element {
	return waitCondition(w, w.we.isVisible, false, "not visible")
}

func (w *waiter) UntilIsEnabled() *Element {
	return waitCondition(w, w.we.isEnabled, true, "enabled")
}

func (w *waiter) UntilIsNotEnabled() *Element {
	return waitCondition(w, w.we.isEnabled, false, "not enabled")
}

func (w *waiter) UntilIsSelected() *Element {
	return waitCondition(w, w.we.isSelected, true, "selected")
}

func (w *waiter) UntilIsNotSelected() *Element {
	return waitCondition(w, w.we.isSelected, false, "not selected")
}

func waitCondition(
	w *waiter,
	condition func() (*Response, error),
	expected bool,
	conditionName string,
) *Element {
	startTime := time.Now()
	endTime := startTime.Add(w.timeout)

	intialSettings := *w.we.settings

	w.we.settings.IgnoreNotFound = true

	defer func() {
		w.we.settings = &intialSettings
	}()

	for {
		if endTime.Before(time.Now()) {
			HandleError(
				errors.Errorf(
					"Element %q is not %s after %s (time elapsed %dms)",
					w.we.Selector,
					conditionName,
					w.timeout,
					time.Since(startTime).Milliseconds(),
				),
			)

			return w.we
		}

		res, err := condition()
		if err != nil {
			if errors.As(err, &ErrWebIDNotSet) {
				time.Sleep(w.we.settings.PollInterval.Duration)

				continue
			}

			HandleError(errors.Wrap(err, "could not get condition"))

			return w.we
		}

		if res.Value.(bool) == expected {
			logger.Infof(
				"Element %q is %s after %s (time elapsed %dms)",
				w.we.Selector,
				conditionName,
				w.timeout,
				time.Since(startTime).Milliseconds(),
			)

			return w.we
		}

		time.Sleep(w.we.settings.PollInterval.Duration)
	}
}

func waitPresent(
	w *waiter,
	bePresent bool,
) *Element {
	startTime := time.Now()
	endTime := startTime.Add(w.timeout)

	conditionName := "present"
	if !bePresent {
		conditionName = "not present"
	}

	for {
		if endTime.Before(time.Now()) {
			HandleError(
				errors.Errorf(
					"Element %q is not %s after %s (time elapsed %s)",
					w.we.Selector,
					conditionName,
					w.timeout,
					time.Since(startTime),
				),
			)

			return w.we
		}

		id, err := w.we.findElement()
		if err != nil {
			time.Sleep(w.we.settings.PollInterval.Duration)

			continue
		}

		if id != "" && bePresent {
			return w.we
		}

		if id == "" && !bePresent {
			return w.we
		}

		time.Sleep(w.we.settings.PollInterval.Duration)
	}
}
