package checks

import (
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/zenoss/zenkit/healthcheck"
)

// FileChecker checks the existence of a file and returns an error
// if the file exists.
func FileChecker(f string) healthcheck.Checker {
	return healthcheck.CheckFunc(func() error {
		if _, err := os.Stat(f); err == nil {
			return errors.New("file exists")
		}
		return nil
	})
}

// HTTPChecker does a HEAD request and verifies that the HTTP status code
// returned matches statusCode.
func HTTPChecker(r string, statusCode int, timeout time.Duration, headers http.Header) healthcheck.Checker {
	return healthcheck.CheckFunc(func() error {
		client := http.Client{
			Timeout: timeout,
		}
		req, err := http.NewRequest("HEAD", r, nil)
		if err != nil {
			return errors.Wrap(err, "error creating request: "+r)
		}
		for headerName, headerValues := range headers {
			for _, headerValue := range headerValues {
				req.Header.Add(headerName, headerValue)
			}
		}
		response, err := client.Do(req)
		if err != nil {
			return errors.Wrap(err, "error while checking: "+r)
		}
		if response.StatusCode != statusCode {
			return errors.New("downstream service returned unexpected status: " + strconv.Itoa(response.StatusCode))
		}
		return nil
	})
}

// TCPChecker attempts to open a TCP connection.
func TCPChecker(addr string, timeout time.Duration) healthcheck.Checker {
	return healthcheck.CheckFunc(func() error {
		conn, err := net.DialTimeout("tcp", addr, timeout)
		if err != nil {
			return errors.Wrap(err, "connection to "+addr+" failed")
		}
		conn.Close()
		return nil
	})
}
