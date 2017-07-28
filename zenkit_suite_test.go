package zenkit_test

import (
	"math/rand"

	"github.com/goadesign/goa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestZenkit(t *testing.T) {
	RegisterFailHandler(Fail)

	rand.Seed(GinkgoRandomSeed())

	RunSpecs(t, "Zenkit Suite")

}

type NullLogAdapter struct{}

func (a *NullLogAdapter) Info(msg string, keyvals ...interface{})   {}
func (a *NullLogAdapter) Error(msg string, keyvals ...interface{})  {}
func (a *NullLogAdapter) New(keyvals ...interface{}) goa.LogAdapter { return a }
