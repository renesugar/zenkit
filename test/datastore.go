package test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/subosito/gotenv"
)

// StartDatastoreEmulator starts the Google Cloud Datastore emulator, running
// on the default port, and sets up environment variables so the SDK will use
// it without configuration. The emulator can be stopped by cancelling the
// context.
func StartDatastoreEmulator(project string) (context.CancelFunc, error) {
	cmd := []string{"gcloud", "beta", "emulators", "datastore", "start",
		"--project", project,
		"--no-store-on-disk",
		"--consistency", "1.0",
		"--host-port", "0.0.0.0:43571",
	}
	ctx, cancel := context.WithCancel(context.Background())
	go exec.CommandContext(ctx, cmd[0], cmd[1:]...).Run()
	envcmd := exec.Command("gcloud", "beta", "emulators", "datastore", "env-init")
	output, err := envcmd.CombinedOutput()
	if err != nil {
		defer cancel()
		return nil, errors.Wrap(err, "unable to run env-init to determine environment")
	}
	if err := gotenv.Apply(bytes.NewReader(output)); err != nil {
		defer cancel()
		return nil, errors.Wrap(err, "unable to apply datastore environment vars")
	}
	for {
		_, err := http.Get(fmt.Sprintf("http://%s/", os.Getenv("DATASTORE_EMULATOR_HOST")))
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		break
	}
	return cancel, nil
}

func ResetDatastore() error {
	addr := os.Getenv("DATASTORE_EMULATOR_HOST")
	_, err := http.Post(fmt.Sprintf("http://%s/reset", addr), "text/plain", strings.NewReader(""))
	return err
}
