package driver

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"github.com/pkg/errors"
)

// Driver resembles a browser driver and parameters to connect to it.
type Driver struct {
	WebdriverPath string
	Port          int
	RemoteURL     string
	cmd           *exec.Cmd
}

func NewDriver(webdriverPath string, port int, remoteURL string) *Driver {
	return &Driver{
		WebdriverPath: webdriverPath,
		Port:          port,
		RemoteURL:     remoteURL,
	}
}

func (d *Driver) Launch() error {
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

	go printLogs(output)

	d.cmd = cmd

	return nil
}

func (d *Driver) Stop() error {
	err := d.cmd.Process.Kill()
	if err != nil {
		return errors.Wrap(err, "failed to kill browser driver")
	}

	return nil
}

func printLogs(output io.ReadCloser) {
	scanner := bufio.NewScanner(output)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
