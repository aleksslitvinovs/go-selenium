package element

import (
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/client"
)

type waiter struct {
	c       *client.Client
	e       *Element
	timeout time.Duration
}

// func (e *Element) WaitUntil(
// 	c *client.Client,
// 	condition func(c *client.Client, e *Element) bool,
// 	timeout time.Duration,
// ) error {
// 	endTime := time.Now().Add(timeout)

// 	for {
// 		if endTime.Before(time.Now()) {
// 			return errors.New("timeout")
// 		}

// 		if condition(c, e) {
// 			return nil
// 		}

// 		// TODO: Implement configurable polling interval
// 		time.Sleep(time.Second)
// 	}
// }

func (e *Element) WaitUntil(c *client.Client, timeout time.Duration) *waiter {
	return &waiter{
		c:       c,
		e:       e,
		timeout: timeout,
	}
}

func (w *waiter) IsVisible() error {
	return w.isVisible(true)
}

func (w *waiter) IsNotVisible() error {
	return w.isVisible(false)
}

func (w *waiter) isVisible(expected bool) error {
	startTime := time.Now()
	endTime := startTime.Add(w.timeout)

	for {
		if endTime.Before(time.Now()) {
			return errors.Errorf(
				// FIXME: Change error message based on expected value
				"element was not visible in %dms",
				time.Since(startTime).Milliseconds(),
			)
		}

		actual, err := w.e.IsVisible(w.c)
		if err != nil {
			return errors.Wrap(err, "could not get visibility state")
		}

		if actual == expected {
			return nil
		}
	}
}

func (w *waiter) IsEnabled() error {
	return w.isEnabled(true)
}

func (w *waiter) IsNotEnabled() error {
	return w.isEnabled(false)
}

func (w *waiter) isEnabled(expected bool) error {
	startTime := time.Now()
	endTime := startTime.Add(w.timeout)

	for {
		if endTime.Before(time.Now()) {
			return errors.Errorf(
				// FIXME: Change error message based on expected value
				"element was not enabled in %dms",
				time.Since(startTime).Milliseconds(),
			)
		}

		actual, err := w.e.IsEnabled(w.c)
		if err != nil {
			return errors.Wrap(err, "could not get enabled state")
		}

		if actual == expected {
			return nil
		}
	}
}

func (w *waiter) IsSelected() error {
	return w.isSelected(true)
}

func (w *waiter) IsNotSelected() error {
	return w.isSelected(false)
}

func (w *waiter) isSelected(expected bool) error {
	startTime := time.Now()
	endTime := startTime.Add(w.timeout)

	for {
		if endTime.Before(time.Now()) {
			return errors.Errorf(
				// FIXME: Change error message based on expected value
				"element was not enabled in %dms",
				time.Since(startTime).Milliseconds(),
			)
		}

		actual, err := w.e.IsSelected(w.c)
		if err != nil {
			return errors.Wrap(err, "could not get selected state")
		}

		if actual == expected {
			return nil
		}
	}
}
