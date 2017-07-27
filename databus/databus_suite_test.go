package databus_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/docker/libcompose/docker"
	lclient "github.com/docker/libcompose/docker/client"
	"github.com/docker/libcompose/docker/container"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/labels"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	proj project.APIProject
	mu   sync.Mutex

	port_map = make(map[string]string)
	logger   = log.New(GinkgoWriter, "[TEST] ", 0)
)

func TestDatabus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Databus Suite")
}

func ResolveAddress(addr string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	if resolved, ok := port_map[addr]; !ok {
		return "", errors.New("unable to resolve address")
	} else {
		return resolved
	}
}

var _ = BeforeSuite(func() {
	var err error
	proj, err = docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{"docker-compose.yml"},
			ProjectName:  "zenkit-databus-test",
		},
	}, nil)
	Ω(err).ShouldNot(HaveOccurred())

	err = proj.Up(context.Background(), options.Up{})
	Ω(err).ShouldNot(HaveOccurred())

	client, err := lclient.Create(lclient.Options{})
	Ω(err).ShouldNot(HaveOccurred())

	containers, err := container.ListByFilter(context.Background(), client, labels.PROJECT.Eq("zenkitdatabustest"))
	Ω(err).ShouldNot(HaveOccurred())

	mu.Lock()
	for _, c := range containers {
		name := c.Labels[string(labels.SERVICE)]
		net := c.NetworkSettings.Networks["bridge"]
		for _, p := range c.Ports {
			port_map[fmt.Sprintf("%s:%d", name, p.PrivatePort)] = fmt.Sprintf("%s:%d", net.IPAddress, p.PrivatePort)
		}
	}
	mu.Unlock()

})

var _ = AfterSuite(func() {
	proj.Down(context.Background(), options.Down{})
})
