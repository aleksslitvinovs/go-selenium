package selenium

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/aleksslitvinovs/go-selenium/logger"
	"github.com/aleksslitvinovs/go-selenium/types"
	"github.com/pkg/errors"
)

// Driver resembles a browser Driver and parameters to connect to it.
type Driver struct {
	webDriverPath string
	remoteURL     string
	port          int
	timeout       *types.Time
	cmd           *exec.Cmd
}

// NewDriver creates a new Driver with the supplied path to the browser driver
// executable and remote URL to use. If remoteURL is empty, ErrInvalidParameters
// error is returned.
func NewDriver(webdriverPath string, remoteURL string) (*Driver, error) {
	if remoteURL == "" {
		return nil, errors.Wrap(
			types.ErrInvalidParameters,
			fmt.Sprintf("remoteURL cannot be %s", remoteURL),
		)
	}

	u, err := url.Parse(remoteURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse remote URL")
	}

	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse port")
	}

	return &Driver{
		webDriverPath: webdriverPath,
		port:          port,
		remoteURL:     remoteURL,
	}, nil
}

// Start starts the browser driver. Driver must be started within the provided
// timeout.
func (d *Driver) Start(timeout *types.Time) error {
	d.timeout = timeout

	if d.webDriverPath == "" {
		return errors.Wrap(
			types.ErrInvalidParameters, "webdriver path cannot be empty",
		)
	}

	if d.port == 0 {
		return errors.Wrap(
			types.ErrInvalidParameters, "port cannot be 0",
		)
	}

	//nolint:gosec
	cmd := exec.Command(d.webDriverPath, fmt.Sprintf("--port=%d", d.port))
	cmd.Stderr = cmd.Stdout

	d.cmd = cmd

	output, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "failed to get stdout pipe")
	}

	err = cmd.Start()
	if err != nil {
		return errors.Wrap(err, "failed to start command")
	}

	ready := make(chan bool, 1)

	go printLogs(ready, d, output)

	select {
	case <-ready:
		return nil
	case <-time.After(d.timeout.Duration):
		return errors.Errorf(
			"failed to start driver within %s", d.timeout.String(),
		)
	}
}

func (d *Driver) stop() error {
	if d.cmd == nil {
		return nil
	}

	err := d.cmd.Process.Kill()
	if err != nil {
		return errors.Wrap(err, "failed to kill browser driver")
	}

	return nil
}

// IsReady returns true if the browser driver is ready to create new sessions.
// An error is returned if there was an issue retrieving driver's status.
// TODO: make public and private methods.
func (d *Driver) IsReady(c *clientParams) (bool, error) {
	var response struct {
		Value struct {
			Ready   bool   `json:"ready"`
			Message string `json:"message"`
		} `json:"value"`
	}

	_, err := c.api.executeRequestCustom(
		http.MethodGet, "/status", struct{}{}, &response,
	)
	if err != nil {
		return false, errors.Wrap(err, "could not get driver status")
	}

	return response.Value.Ready, nil
}

func printLogs(ready chan<- bool, d *Driver, output io.ReadCloser) {
	scanner := bufio.NewScanner(output)

	for scanner.Scan() {
		line := scanner.Text()

		fmt.Println(line)

		// TODO: Improve error handling
		if strings.Contains(line, "Address already in use") {
			logger.Errorf(
				"Cannot start browser driver. Port %d is already in use.",
				d.port,
			)

			d.stop() //nolint:errcheck
		}

		// Chromedriver is ready
		if strings.Contains(line, "ChromeDriver was started successfully") ||
			// GeckoDriver is ready
			strings.Contains(line, "Listening on") {
			ready <- true
		}
	}
}
