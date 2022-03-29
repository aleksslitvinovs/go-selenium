package types

type Handle struct {
	ID   string `json:"handle"`
	Type string `json:"type"`
}

type Sessioner interface {
	NewElement(selector string) WebElementer

	// OpenURL opens a new window with the given URL.
	OpenURL(url string)

	// GetCurrentURL returns the current URL of the browsing context.
	GetCurrentURL() string

	// Refresh refreshes the current page.
	Refresh()

	// Back navigates back in the browser history.
	Back()

	// Forward navigates forward in the browser history.
	Forward()

	// GetTitle returns the current page title.
	GetTitle() string

	// GetWindowHandle returns the current browsing context window handle.
	GetWindowHandle() string

	// GetWindowHandles returns an array of window handles for each browsing
	// context. The order of the handles is arbitrary.
	GetWindowHandles() []string

	// CloseWindow closes the current browsing context.
	CloseWindow()

	// SwitchHandle switches to the given browsing context.
	SwitchHandle(handle string)

	// NewTab opens a new tab in the current window.
	NewTab() *Handle

	// NewWindow opens a window.
	NewWindow() *Handle

	// TODO: Not working yet.
	// SwitchToFrame switches to the given iframe element.
	SwitchToFrame(id int)

	// SwitchToParentFrame switches to the parent frame of the current frame.
	SwitchToParentFrame()

	DeleteSession()

	AddError(err string)
	GetID() string
}
