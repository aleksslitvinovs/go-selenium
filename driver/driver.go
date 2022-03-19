package driver

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"github.com/pkg/errors"
)

type Driver struct {
	WebdriverPath string
	Port          int
	RemoteURL     string
	cmd           *exec.Cmd
}

// NewDriverBuilder creates a new driver builder with default configuration.
func NewDriverBuilder() *Driver {
	return &Driver{Port: 4444, RemoteURL: "http://localhost"}
}

func (d *Driver) SetDriver(path string) *Driver {
	d.WebdriverPath = path
	return d
}

func (d *Driver) SetPort(port int) *Driver {
	d.Port = port
	return d
}

func (d *Driver) SetRemoteUrl(url string) *Driver {
	d.RemoteURL = url
	return d
}

func (d *Driver) Build() *Driver {
	return d
}

func (d *Driver) Launch() error {
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

func (c *Driver) Stop() error {
	err := c.cmd.Process.Kill()
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
