package databus_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDatabus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Databus Suite")
}
