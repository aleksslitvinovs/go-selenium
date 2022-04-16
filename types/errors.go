package types

import "errors"

var (
	// WebDriver errors.

	ErrElementlickIntercepted = errors.New("element click intercepted")
	ErrElementNotInteractable = errors.New("element not interactable")
	ErrInsecureCertificate    = errors.New("insecure certificate")
	ErrInvalidArgument        = errors.New("invalid argument")
	ErrInvalidCookieDomain    = errors.New("invalid cookie domain")
	ErrInvalidElementState    = errors.New("invalid element state")
	ErrInvalidSelector        = errors.New("invalid selector")
	ErrInvalidSessionID       = errors.New("invalid session id")
	ErrJavaScriptError        = errors.New("javascript error")
	ErrMoveTargetOutOfBounds  = errors.New("move target out of bounds")
	ErrNoSuchAlert            = errors.New("no such alert")
	ErrNoSuchCookie           = errors.New("no such cookie")
	ErrNoSuchElement          = errors.New("no such element")
	ErrNoSuchFrame            = errors.New("no such frame")
	ErrNoSuchWindow           = errors.New("no such window")
	ErrScriptTimeout          = errors.New("script timeout")
	ErrSessionNotCreated      = errors.New("session not created")
	ErrStaleElementReference  = errors.New("stale element reference")
	ErrDetachedShadowRoot     = errors.New("detached shadow root")
	ErrTimeout                = errors.New("timeout")
	ErrUnableToSetCookie      = errors.New("unable to set cookie")
	ErrUnableToCaptureScreen  = errors.New("unable to capture screen")
	ErrUnexpectedAlertOpen    = errors.New("unexpected alert open")
	ErrUnknownCommand         = errors.New("unknown command")
	ErrUnknownError           = errors.New("unknown error")
	ErrUnknownMethod          = errors.New("unknown method")
	ErrUnsupportedOperation   = errors.New("unsupported operation")

	// Go-Selenium errors.

	ErrInvalidParameters = errors.New("invalid parameters")
	ErrFailedRequest     = errors.New("failed to execute request")
)
