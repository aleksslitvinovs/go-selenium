package selenium

import (
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/logger"
)

type waiter struct {
	e       *element
	timeout time.Duration
}

func (e *element) WaitFor(timeout time.Duration) *waiter {
	return &waiter{
		e:       e,
		timeout: timeout,
	}
}

func (w *waiter) UntilIsPresent() *element {
	return waitPresent(w, true)
}
func (w *waiter) UntilIsNotPresent() *element {
	return waitPresent(w, false)
}

func (w *waiter) UntilIsVisible() *element {
	return waitCondition(w, w.e.isVisible, true, "visible")
}
func (w *waiter) UntilIsNotVisible() *element {
	return waitCondition(w, w.e.isVisible, false, "not visible")
}

func (w *waiter) UntilIsEnabled() *element {
	return waitCondition(w, w.e.isEnabled, true, "enabled")
}

func (w *waiter) UntilIsNotEnabled() *element {
	return waitCondition(w, w.e.isEnabled, false, "not enabled")
}

func (w *waiter) UntilIsSelected() *element {
	return waitCondition(w, w.e.isSelected, true, "selected")
}

func (w *waiter) UntilIsNotSelected() *element {
	return waitCondition(w, w.e.isSelected, false, "not selected")
}

func waitCondition(
	w *waiter,
	condition func() (*response, error),
	expected bool,
	conditionName string,
) *element {
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

		if res.Value.(bool) == expected {
			logger.Infof(
				"Element %q is %s after %s (time elapsed %dms)",
				w.e.Selector,
				conditionName,
				w.timeout,
				time.Since(startTime).Milliseconds(),
			)

			return w.e
		}

		time.Sleep(w.e.settings.PollInterval.Duration)
	}
}

func waitPresent(
	w *waiter,
	bePresent bool,
) *element {
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
			return w.e
		}

		time.Sleep(w.e.settings.PollInterval.Duration)
	}
}
