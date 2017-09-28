package funcname_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFuncname(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Funcname Suite")
}
