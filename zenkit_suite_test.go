package zenkit_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestZenkit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Zenkit Suite")
}
