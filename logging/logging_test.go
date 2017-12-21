package logging_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"net/http"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	. "github.com/zenoss/zenkit/logging"
	"github.com/zenoss/zenkit/test"

	"encoding/json"
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

	Context("with the service logger with request id", func() {
		var svc *goa.Service

		BeforeEach(func() {
			svc = goa.New(test.RandString(8))
			svc.Use(middleware.RequestID())
			svc.WithLogger(ServiceLogger())
		})

		It("should produce a logger with logrus entry containing 'req_id'", func() {
			var logger *logrus.Entry
			entry := ContextLoggerWithReqId(svc.Context)
			Ω(entry).ShouldNot(BeNil())
			Ω(entry).Should(BeAssignableToTypeOf(logger))
			Ω(entry.Data).Should(HaveLen(1))
			Ω(entry.Data).Should(HaveKey("req_id"))
			Ω(entry.Data["req_id"]).Should(Not(BeNil()))
		})
	})

	Context("with the error-response middleware logger", func() {
		var (
			svc *goa.Service
			resp   http.ResponseWriter
			req    *http.Request
			b      bytes.Buffer
			logger *logrus.Logger
		)

		okHandler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			goa.ContextResponse(ctx).WriteHeader(200)
			goa.ContextResponse(ctx).Write([]byte{})
			return nil
		}

		BeforeEach(func() {
			svc = goa.New(test.RandString(8))
			svc.WithLogger(ServiceLogger())
			logger = ContextLogger(svc.Context).Logger
			resp   = httptest.NewRecorder()
			req, _ = http.NewRequest("", "http://example.com", nil)
		})

		It("shouldn't panic if there's no logger defined", func() {
			svc.WithLogger(nil)
			l := LogErrorResponse()(okHandler)
			ctx := goa.NewContext(svc.Context, resp, req, nil)

			err := l(ctx, resp, req)

			Ω(err).ShouldNot(HaveOccurred())
		})

		It("shouldn't log anything if level != debug", func() {
			SetLogLevel(svc, "info")
			logger.Out = &b
			l := LogErrorResponse()(okHandler)
			ctx := goa.NewContext(svc.Context, resp, req, nil)

			err := l(ctx, resp, req)

			Ω(err).ShouldNot(HaveOccurred())
			s := b.String()
			Ω(s).Should(BeEmpty())
		})

		It("shouldn't log anything if the response is OK", func() {
			SetLogLevel(svc, "debug")
			logger.Out = &b
			l := LogErrorResponse()(okHandler)
			ctx := goa.NewContext(svc.Context, resp, req, nil)

			err := l(ctx, resp, req)

			Ω(err).ShouldNot(HaveOccurred())
			s := b.String()
			Ω(s).Should(BeEmpty())
		})

		It("should log something if the response is !OK", func() {
			var errResponse *goa.ErrorResponse
			SetLogLevel(svc, "debug")
			logger.Out = &b
			failHandler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
				// Simulate GOA detecting a bad request
				var err error
				err = goa.ErrBadRequest("fail test")
				errResponse = err.(*goa.ErrorResponse)
				ctx = goa.WithError(ctx, errResponse)
				resp := goa.ContextResponse(ctx)
				resp.ErrorCode = errResponse.Token()
				rw.WriteHeader(400)
				body, _ := json.Marshal(goa.ContextError(ctx))
				rw.Write(body)
				return err
			}
			l := LogErrorResponse()(failHandler)
			ctx := goa.NewContext(svc.Context, resp, req, nil)

			err := l(ctx, goa.ContextResponse(ctx), req)

			Ω(err).Should(HaveOccurred())
			s := b.String()
			Ω(s).Should(Not(BeEmpty()))
			Ω(s).Should(ContainSubstring("returned an error"))
			Ω(s).Should(ContainSubstring(fmt.Sprintf("req_id=%s", middleware.ContextRequestID(ctx))))
			Ω(s).Should(ContainSubstring(fmt.Sprintf("error_id=%s", errResponse.ID)))
			Ω(s).Should(ContainSubstring(fmt.Sprintf("code=%s", errResponse.Code)))
			Ω(s).Should(ContainSubstring(fmt.Sprintf("status=%d", errResponse.Status)))
			Ω(s).Should(ContainSubstring(fmt.Sprintf("detail=\"%s\"", errResponse.Detail)))
		})

		It("should log something if the response is not an instance of ErrorResponse", func() {
			var errResponse *goa.ErrorResponse
			SetLogLevel(svc, "debug")
			logger.Out = &b
			failHandler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
				// Simulate GOA detecting a bad request
				var err error
				err = goa.ErrBadRequest("fail test")
				errResponse = err.(*goa.ErrorResponse)
				ctx = goa.WithError(ctx, errResponse)
				resp := goa.ContextResponse(ctx)
				resp.ErrorCode = errResponse.Token()
				rw.WriteHeader(400)

				// Do NOT write ErrorResponse to the response body
				rw.Write([]byte{})
				return err
			}
			l := LogErrorResponse()(failHandler)
			ctx := goa.NewContext(svc.Context, resp, req, nil)

			err := l(ctx, goa.ContextResponse(ctx), req)

			Ω(err).Should(HaveOccurred())
			s := b.String()
			Ω(s).Should(Not(BeEmpty()))
			Ω(s).Should(ContainSubstring("Unable to unmarshall buffer into ErrorResponse"))
			Ω(s).Should(ContainSubstring(fmt.Sprintf("req_id=%s", middleware.ContextRequestID(ctx))))
		})

	})

})
