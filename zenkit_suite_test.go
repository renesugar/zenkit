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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type NullLogAdapter struct{}

func (a *NullLogAdapter) Info(msg string, keyvals ...interface{})   {}
func (a *NullLogAdapter) Error(msg string, keyvals ...interface{})  {}
func (a *NullLogAdapter) New(keyvals ...interface{}) goa.LogAdapter { return a }
