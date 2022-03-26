package driver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/logger"

	"github.com/pkg/errors"
)

type Opts struct {
	Timeout time.Duration
}

// Driver resembles a browser driver and parameters to connect to it.
type Driver struct {
	WebdriverPath string
	Port          int
	RemoteURL     string
	Opts          *Opts
	cmd           *exec.Cmd
}

func NewDriver(webdriverPath string, port int, remoteURL string, opts *Opts) *Driver {
	return &Driver{
		WebdriverPath: webdriverPath,
		Port:          port,
		RemoteURL:     remoteURL,
		Opts:          opts,
	}
}

func (d *Driver) Start() error {
	//nolint:gosec
	cmd := exec.Command(d.WebdriverPath, fmt.Sprintf("--port=%d", d.Port))

	output, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "failed to get stdout pipe")
	}

	cmd.Stderr = cmd.Stdout

	err = cmd.Start()
	if err != nil {
		return errors.Wrap(err, "failed to start command")
	}

	ready := make(chan bool, 1)

	go printLogs(ready, d, output)

	d.cmd = cmd

	select {
	case <-ready:
	case <-time.After(d.Opts.Timeout):
		return errors.Errorf("failed to start driver within %s", d.Opts.Timeout)
	}

	return nil
}

func (d *Driver) Stop() error {
	err := d.cmd.Process.Kill()
	if err != nil {
		return errors.Wrap(err, "failed to kill browser driver")
	}

	return nil
}

func IsReady(c api.Requester) (bool, error) {
	res, err := api.ExecuteRequestRaw(
		http.MethodGet, "/status", c, struct{}{},
	)
	if err != nil {
		return false, errors.Wrap(err, "failed to get status")
	}

	type response struct {
		Value struct {
			Ready   bool   `json:"ready"`
			Message string `json:"message"`
		} `json:"value"`
	}

	var r response

	if err := json.Unmarshal(res, &r); err != nil {
		return false, errors.Wrap(err, "failed to unmarshal response")
	}

	return r.Value.Ready, nil
}
func printLogs(ready chan<- bool, d *Driver, output io.ReadCloser) {
	scanner := bufio.NewScanner(output)

	for scanner.Scan() {
		line := scanner.Text()

		fmt.Println(line)

		// TODO: Improve error handling
		if strings.Contains(line, "Address already in use") {
			logger.Fatalf(
				"Cannot start browser driver. Port %d is already in use.",
				d.Port,
			)

			d.Stop()
		}

		// TODO: Add handling for FF
		if strings.Contains(line, "ChromeDriver was started successfully") {
			ready <- true
		}
	}
}
