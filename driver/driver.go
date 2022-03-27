package driver

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/api"
	"github.com/theRealAlpaca/go-selenium/config"
	"github.com/theRealAlpaca/go-selenium/logger"
)

// Driver resembles a browser driver and parameters to connect to it.
type Driver struct {
	WebDriverPath string
	Port          int
	RemoteURL     string
	Timeout       time.Duration
	cmd           *exec.Cmd
}

func newDriver(
	webdriverPath string, port int, remoteURL string, timeout time.Duration,
) *Driver {
	return &Driver{
		WebDriverPath: webdriverPath,
		Port:          port,
		RemoteURL:     remoteURL,
		Timeout:       timeout,
	}
}

func Start(conf *config.WebDriverConfig) (*Driver, error) {
	d := newDriver(conf.PathToBinary, conf.Port, conf.URL, conf.Timeout)

	if d.Port == 0 {
		d.Port = 4444
	}

	if d.RemoteURL == "" {
		d.RemoteURL = "http://localhost"
	}

	if d.Timeout == 0 {
		d.Timeout = time.Second * 10
	}

	//nolint:gosec
	cmd := exec.Command(d.WebDriverPath, fmt.Sprintf("--port=%d", d.Port))
	cmd.Stderr = cmd.Stdout

	d.cmd = cmd

	output, err := cmd.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stdout pipe")
	}

	err = cmd.Start()
	if err != nil {
		return nil, errors.Wrap(err, "failed to start command")
	}

	ready := make(chan bool, 1)

	go printLogs(ready, d, output)

	select {
	case <-ready:
	case <-time.After(d.Timeout):
		return nil, errors.Errorf("failed to start driver within %s", d.Timeout)
	}

	return d, nil
}

func (d *Driver) Stop() error {
	err := d.cmd.Process.Kill()
	if err != nil {
		return errors.Wrap(err, "failed to kill browser driver")
	}

	return nil
}

func IsReady(c api.Requester) (bool, error) {
	var response struct {
		Value struct {
			Ready   bool   `json:"ready"`
			Message string `json:"message"`
		} `json:"value"`
	}

	err := api.ExecuteRequestCustom(
		http.MethodGet, "/status", c, struct{}{}, &response,
	)
	if err != nil {
		return false, errors.Wrap(err, "failed to get status")
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
			logger.Fatalf(
				"Cannot start browser driver. Port %d is already in use.",
				d.Port,
			)

			d.Stop() //nolint:errcheck
		}

		// TODO: Add handling for FF
		if strings.Contains(line, "ChromeDriver was started successfully") {
			ready <- true
		}
	}
}
