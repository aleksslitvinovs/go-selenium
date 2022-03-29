package webelement

import (
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/types"
	"github.com/theRealAlpaca/go-selenium/util"
)

type waiter struct {
	we      *webElement
	timeout time.Duration
}

var (
	ErrWebIDNotSet = errors.New("WebID not set")

	_ (types.Waiterer) = (*waiter)(nil)
)

func (we *webElement) WaitFor(timeout time.Duration) types.Waiterer {
	return &waiter{
		we:      we,
		timeout: timeout,
	}
}

func (w *waiter) UntilIsVisible() types.WebElementer {
	return waitCondition(w, w.we.isVisible, true, "is visible")
}
func (w *waiter) UntilIsNotVisible() types.WebElementer {
	return waitCondition(w, w.we.isVisible, false, "not visible")
}

func (w *waiter) UntilIsEnabled() types.WebElementer {
	return waitCondition(w, w.we.isEnabled, true, "enabled")
}

func (w *waiter) UntilIsNotEnabled() types.WebElementer {
	return waitCondition(w, w.we.isEnabled, false, "not enabled")
}

func (w *waiter) UntilIsSelected() types.WebElementer {
	return waitCondition(w, w.we.isSelected, true, "selected")
}

func (w *waiter) UntilIsNotSelected() types.WebElementer {
	return waitCondition(w, w.we.isSelected, false, "not selected")
}

func waitCondition(
	w *waiter,
	condition func() (*api.Response, error),
	expected bool,
	conditionName string,
) types.WebElementer {
	startTime := time.Now()
	endTime := startTime.Add(w.timeout)

	intialSettings := *w.we.settings

	w.we.settings.IgnoreNotFound = true

	defer func() {
		w.we.settings = &intialSettings
	}()

	for {
		if endTime.Before(time.Now()) {
			util.HandleError(
				w.we.session,
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

		res, err := condition()
		if err != nil {
			if errors.As(err, &ErrWebIDNotSet) {
				time.Sleep(w.we.settings.PollInterval.Duration)

				continue
			}

			util.HandleError(
				w.we.session, errors.Wrap(err, "could not get condition"),
			)

			return w.we
		}

		if res.Value.(bool) == expected {
			return w.we
		}

		time.Sleep(w.we.settings.PollInterval.Duration)
	}
}
