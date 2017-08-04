package zenkit_test

import (
	"bytes"
	"context"
	"fmt"

	"github.com/goadesign/goa"
	"github.com/sirupsen/logrus"
	. "github.com/zenoss/zenkit"
	"github.com/zenoss/zenkit/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TheTestFunction(ctx context.Context) {
	defer LogEntryAndExit(ctx)()
	logger := ContextLogger(ctx)
	if logger != nil {
		logger.Info("Inside!")
	}
}

var _ = Describe("Logging", func() {

	Context("with the service logger", func() {

		var svc *goa.Service

		BeforeEach(func() {
			svc = goa.New(test.RandString(8))
			svc.WithLogger(ServiceLogger())
		})

		It("should produce a logger that can be used as a logrus entry", func() {
			var logger *logrus.Entry
			entry := ContextLogger(svc.Context)
			Ω(entry).ShouldNot(BeNil())
			Ω(entry).Should(BeAssignableToTypeOf(logger))
		})

		It("should be able to set the log level to a valid level", func() {
			logger := ContextLogger(svc.Context).Logger
			var b bytes.Buffer
			logger.Out = &b

			SetLogLevel(svc, "error")
			s := b.String()
			Ω(s).Should(BeEmpty()) // Because the new log level is higher than that at which we log level changes

			b.Reset()
			SetLogLevel(svc, "info")
			s = b.String()
			Ω(s).Should(ContainSubstring("Log level changed"))
			Ω(s).Should(ContainSubstring("newlevel=info"))

			SetLogLevel(svc, "error")

			b.Reset()
			msg := test.RandString(8)
			svc.LogInfo(msg)
			s = b.String()
			Ω(s).ShouldNot(ContainSubstring(msg))

			b.Reset()
			msg = test.RandString(8)
			svc.LogError(msg)
			s = b.String()
			Ω(s).Should(ContainSubstring(msg))
		})

		It("should fail to set the log level to an invalid level", func() {
			logger := ContextLogger(svc.Context).Logger
			var b bytes.Buffer
			logger.Out = &b

			level := test.RandString(8)
			SetLogLevel(svc, level)
			s := b.String()
			Ω(s).Should(ContainSubstring("Unable to parse log level. Not changing."))
			Ω(s).Should(ContainSubstring(fmt.Sprintf("badlevel=%s", level)))
		})

		It("should decline to reset the log level to the existing log level", func() {
			logger := ContextLogger(svc.Context).Logger
			var b bytes.Buffer
			logger.Out = &b
			SetLogLevel(svc, "debug")
			b.Reset()

			SetLogLevel(svc, "debug")
			s := b.String()
			Ω(s).Should(ContainSubstring("Requested log level is already active. Ignoring."))
		})

		It("should trace log appropriately", func() {
			logger := ContextLogger(svc.Context).Logger
			var b bytes.Buffer
			logger.Out = &b
			SetLogLevel(svc, "debug")
			b.Reset()

			TheTestFunction(svc.Context)
			s := b.String()
			Ω(s).Should(ContainSubstring("ENTER TheTestFunction()"))
			Ω(s).Should(ContainSubstring("Inside!"))
			Ω(s).Should(ContainSubstring("EXIT TheTestFunction()"))
		})

		It("shouldn't panic the trace logger if there's no logger defined", func() {
			TheTestFunction(context.Background())
		})
	})

})
