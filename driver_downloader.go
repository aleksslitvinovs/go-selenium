package selenium

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/theRealAlpaca/go-selenium/logger"
)

const (
	chromedriver = "chromedriver"
	geckodriver  = "geckodriver"
)

var (
	//nolint:revive
	ErrUnsupportedPlatform = errors.Errorf(
		"unsupported platform %s %s", runtime.GOOS, runtime.GOARCH,
	)
	chromeDriverURL  = "https://chromedriver.storage.googleapis.com/%s/chromedriver_%s"                          //nolint:lll
	firefoxDriverURL = "https://github.com/mozilla/geckodriver/releases/download/v0.30.0/geckodriver-v0.30.0-%s" //nolint:lll
	defaultPath      = "./"
)

func downloadDriver(driverName string) error {
	// Default driver is found on the system.
	if _, err := os.Stat(path.Join(defaultPath, driverName)); err == nil {
		return nil
	}

	var platform string

	var binaryData []byte

	switch driverName {
	case chromedriver:
		platform, err := getPlatformForChrome()
		if err != nil {
			return errors.Wrap(err, "failed to get platform")
		}

		binaryData, err = downloadChrome(platform)
		if err != nil {
			return errors.Wrap(err, "failed to download chromedriver")
		}
	case geckodriver:
		platform, err := getPlatformForFirefox()
		if err != nil {
			return errors.Wrap(err, "failed to get platform")
		}

		binaryData, err = downloadFirefox(platform)
		if err != nil {
			return errors.Wrap(err, "failed to download geckodriver")
		}
	default:
		return errors.Errorf("unsupported driver: %s", driverName)
	}

	err := saveBinary(binaryData, platform, driverName)
	if err != nil {
		return errors.Wrap(err, "failed to save binary")
	}

	return nil
}

func downloadChrome(platform string) ([]byte, error) {
	logger.Info("Downloading Chrome driver")

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"https://chromedriver.storage.googleapis.com/LATEST_RELEASE",
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(
			err,
			"could not get create request to get latest chromedriver version",
		)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(
			err,
			"could not get latest chromedriver version",
		)
	}
	defer res.Body.Close()

	version, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read latest release version")
	}

	req, err = http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		fmt.Sprintf(chromeDriverURL, string(version), platform),
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(
			err, "failed to create request to download chromedriver",
		)
	}

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(
			err, "failed to execute request to download chromedriver",
		)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to download chromedriver")
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read chromedriver binary data")
	}

	return b, nil
}

func downloadFirefox(platform string) ([]byte, error) {
	logger.Info("Downloading Firefox driver")

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		fmt.Sprintf(firefoxDriverURL, platform),
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(
			err, "failed to create request to get geckodriver",
		)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(
			err, "failed to execute request to get geckodriver",
		)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to download geckodriver")
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read geckodriver binary data")
	}

	return b, nil
}

func saveBinary(
	data []byte, platform, driverName string,
) error {
	archiveName := fmt.Sprintf("%s.zip", driverName)

	if driverName == geckodriver {
		if strings.Contains(platform, "win") {
			archiveName = fmt.Sprintf("%s.zip", driverName)
		} else {
			archiveName = fmt.Sprintf("%s.tar.gz", driverName)
		}
	}

	f, err := os.Create(archiveName)
	if err != nil {
		return errors.Wrap(err, "failed to create browser driver file")
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return errors.Wrap(err, "failed to copy browser driver data to a file")
	}

	err = unarchive(archiveName, driverName)
	if err != nil {
		return errors.Wrap(err, "failed to unarchive browser driver")
	}

	err = os.Remove(archiveName)
	if err != nil {
		return errors.Wrap(err, "failed to delete Apimation zip archive")
	}

	return nil
}

func unarchive(archiveName, driverName string) error {
	f, err := os.OpenFile(
		driverName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0766,
	)
	if err != nil {
		return errors.Wrap(err, "failed to open browser driver target file")
	}
	defer f.Close()

	if strings.Contains(archiveName, "zip") {
		err = unzip(f, archiveName)
		if err != nil {
			return errors.Wrap(err, "failed to unzip browser driver")
		}
	}

	if strings.Contains(archiveName, "tar") {
		err := untar(f, archiveName)
		if err != nil {
			return errors.Wrap(err, "failed to untar browser driver")
		}
	}

	return nil
}

func unzip(targetFile *os.File, zipName string) error {
	r, err := zip.OpenReader(zipName)
	if err != nil {
		return errors.Wrap(err, "failed to open browser driver zip file")
	}
	defer r.Close()

	if len(r.File) != 1 {
		return errors.Errorf("Expected 1 file, found %d", len(r.File))
	}

	for _, f := range r.File {
		fr, err := f.Open()
		if err != nil {
			return errors.Wrap(err, "failed to open browser driver file")
		}

		for {
			_, err = io.CopyN(targetFile, fr, 1024)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				return errors.Wrap(err, "failed to copy browser driver file")
			}
		}
	}

	return nil
}

func untar(targetFile *os.File, zipName string) error {
	f, err := os.Open(zipName)
	if err != nil {
		return errors.Wrap(err, "failed to open browser driver file")
	}

	r, err := gzip.NewReader(f)
	if err != nil {
		return errors.Wrap(err, "failed to create tar reader")
	}
	defer r.Close()

	tr := tar.NewReader(r)

	for {
		h, err := tr.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return errors.Wrap(err, "failed to read tar file")
		}

		if h.Typeflag != tar.TypeReg {
			continue
		}

		for {
			_, err = io.CopyN(targetFile, tr, 1024)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				return errors.Wrap(err, "failed to copy browser driver file")
			}
		}
	}

	return nil
}

func getPlatformForChrome() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		if runtime.GOARCH == "amd64" {
			return "mac64.zip", nil
		}

		if runtime.GOARCH == "arm64" {
			return "macm64_m1.zip", nil
		}

		return "", ErrUnsupportedPlatform
	case "linux":
		if strconv.IntSize == 64 {
			return "linux64.zip", ErrUnsupportedPlatform
		}

		return "", ErrUnsupportedPlatform
	case "windows":
		if strconv.IntSize == 32 {
			return "win32.zip", nil
		}

		return "", ErrUnsupportedPlatform
	default:
		return "", ErrUnsupportedPlatform
	}
}

//nolint:cyclop
func getPlatformForFirefox() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		if runtime.GOARCH == "amd64" {
			return "macos.tar.gz", nil
		}

		if runtime.GOARCH == "arm64" {
			return "macos-aarch64.tar.gz", nil
		}

		return "", ErrUnsupportedPlatform
	case "linux":
		if strconv.IntSize == 64 {
			return "linux64.tar.gz", nil
		}

		if strconv.IntSize == 32 {
			return "linux32.tar.gz", nil
		}

		return "", ErrUnsupportedPlatform
	case "windows":
		if strconv.IntSize == 64 {
			return "win64.zip", nil
		}

		if strconv.IntSize == 32 {
			return "win32.zip", nil
		}

		return "", ErrUnsupportedPlatform
	default:
		return "", ErrUnsupportedPlatform
	}
}
