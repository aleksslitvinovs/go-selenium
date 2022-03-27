package util

import (
	"fmt"
	"strings"
	"sync"
)

func DownloadDriver(wg *sync.WaitGroup, name string) string {
	var driverName string

	switch strings.ToLower(name) {
	case "chrome":
		driverName = "chromedriver"
	case "firefox":
		driverName = "geckodriver"
	default:
		panic("unknown driver name")
	}

	wg.Add(1)

	defer wg.Done()

	fmt.Println(driverName)

	panic("not implemented")
}
