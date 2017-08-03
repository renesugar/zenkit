package test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/docker/libcompose/config"
	"github.com/docker/libcompose/docker"
	lclient "github.com/docker/libcompose/docker/client"
	"github.com/docker/libcompose/docker/container"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/labels"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
	"github.com/pkg/errors"
)

var (
	// nameRegexp is a copy of a pattern from libcompose, because the sanitized
	// name isn't gettable from the project for some reason, so we have to run
	// it ourselves
	nameRegexp = regexp.MustCompile("[^a-z0-9]+")
	// ErrNoContainerFound is raised when we don't find any containers using the filters specified
	ErrNoContainerFound = errors.New("no containers found")
)

type Harness interface {
	Start() error
	Stop() error
	Wait(healthcheck func() error, timeout time.Duration) error
	Resolve(service string, port uint64) (string, error)
}

type dockerComposeHarness struct {
	name    string
	project project.APIProject
}

func normalizeName(name string) string {
	return nameRegexp.ReplaceAllString(strings.ToLower(name), "")
}

func NewDockerComposeHarness(name string, dockerComposeFiles ...string) (Harness, error) {
	name = normalizeName(name)
	proj, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: dockerComposeFiles,
			ProjectName:  name,
		},
	}, &config.ParseOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "unable to create libcompose project")
	}
	return &dockerComposeHarness{
		name:    name,
		project: proj,
	}, nil
}

func (h *dockerComposeHarness) Start() error {
	if err := h.project.Up(context.Background(), options.Up{}); err != nil {
		return errors.Wrap(err, "unable to bring up the libcompose project")
	}
	return nil
}

func (h *dockerComposeHarness) Stop() error {
	if err := h.project.Down(context.Background(), options.Down{}); err != nil {
		return errors.Wrap(err, "unable to shut down the libcompose project")
	}
	return nil
}

func (h *dockerComposeHarness) Resolve(service string, port uint64) (string, error) {
	client, err := lclient.Create(lclient.Options{})
	if err != nil {
		return "", errors.Wrap(err, "unable to create Docker client")
	}
	filter := labels.And(labels.PROJECT.Eq(h.name), labels.SERVICE.Eq(service))
	containers, err := container.ListByFilter(context.Background(), client, filter)
	if err != nil {
		return "", errors.Wrap(err, "unable to filter service containers")
	}
	if len(containers) == 0 {
		return "", errors.WithStack(ErrNoContainerFound)
	}
	net := containers[0].NetworkSettings.Networks["bridge"]
	return fmt.Sprintf("%s:%d", net.IPAddress, port), nil
}

func (h *dockerComposeHarness) Wait(healthcheck func() error, timeout time.Duration) error {
	boff := backoff.NewExponentialBackOff()
	boff.MaxElapsedTime = timeout
	boff.InitialInterval = 500 * time.Millisecond
	boff.Multiplier = 1.0
	if err := backoff.Retry(healthcheck, boff); err != nil {
		return errors.Wrap(err, "health check didn't pass within the timeout")
	}
	return nil
}
