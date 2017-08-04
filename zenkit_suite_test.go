package zenkit_test

import (
	"math/rand"

	"github.com/goadesign/goa"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"

	"testing"
)

func TestZenkit(t *testing.T) {
	RegisterFailHandler(Fail)

	rand.Seed(GinkgoRandomSeed())
	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Zenkit suite", []Reporter{junitReporter})
}

type NullLogAdapter struct{}

func (a *NullLogAdapter) Info(msg string, keyvals ...interface{})   {}
func (a *NullLogAdapter) Error(msg string, keyvals ...interface{})  {}
func (a *NullLogAdapter) New(keyvals ...interface{}) goa.LogAdapter { return a }
