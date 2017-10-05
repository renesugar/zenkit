package databus

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/zenoss/zenkit/healthcheck"
)

// KafkaChecker verifies a connection to a kafka broker
func KafkaChecker(addrs ...string) healthcheck.Checker {
	return healthcheck.CheckFunc(func() error {
		client, err := sarama.NewClient(addrs, nil)
		if err != nil {
			return errors.WithStack(err)
		}
		client.Close()
		return nil
	})
}

// SchemaRegistryChecker does a GET request and verifies that the response is
// valid.
func SchemaRegistryChecker(addr string, timeout time.Duration) healthcheck.Checker {
	return healthcheck.CheckFunc(func() error {
		client := http.Client{
			Timeout: timeout,
		}
		r := fmt.Sprintf("http://%s/config", addr)
		req, err := http.NewRequest("GET", r, nil)
		if err != nil {
			return errors.Wrapf(err, "error creating request: "+r)
		}
		response, err := client.Do(req)
		if err != nil {
			return errors.Wrap(err, "error while checking: "+r)
		}
		defer response.Body.Close()
		data, _ := ioutil.ReadAll(response.Body)
		if bytes.Compare(data, []byte(`{"compatibilityLevel":"BACKWARD"}`)) != 0 {
			return errors.New("schema registry is not ok")
		}
		return nil
	})
}
