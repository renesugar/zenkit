package databus_test

import (
	"context"
	"log"
	"time"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	proj   project.APIProject
	ctx_go = context.Background()
	logger = log.New(GinkgoWriter, "[TEST] ", 0)
)

func TestDatabus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Databus Suite")
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

	err = proj.Up(ctx_go, options.Up{})
	Ω(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	time.Sleep(10 * time.Second)
	proj.Down(ctx_go, options.Down{})
})
