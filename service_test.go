package zenkit_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/goatest"
	"github.com/goadesign/goa/middleware"
	. "github.com/zenoss/zenkit"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {

	var (
		svc        *goa.Service
		req        *http.Request
		rw         *httptest.ResponseRecorder
		name       string
		ctrl       *goa.Controller
		resp       interface{}
		logBuf     bytes.Buffer
		respSetter goatest.ResponseSetterFunc = func(r interface{}) { resp = r }
		encoder                               = func(io.Writer) goa.Encoder { return respSetter }
		logger                                = log.New(&logBuf, "", log.Ltime)
	)

	BeforeEach(func() {
		name = RandStringRunes(8)
		req, _ = http.NewRequest("", "http://example.com/", nil)
		rw = httptest.NewRecorder()
	})

	JustBeforeEach(func() {
		svc.WithLogger(goa.NewLogger(logger))
		svc.Encoder = goa.NewHTTPEncoder()
		svc.Encoder.Register(encoder, "*/*")
		ctrl = svc.NewController("test")
	})

	AfterEach(func() {
		logBuf.Reset()
	})

	type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

	RunHandler := func(h HandlerFunc) {
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			defer GinkgoRecover()
			h(ctx, rw, req)
			return nil
		}
		ctrl.MuxHandler("test", handler, nil)(rw, req, url.Values{})
	}

	Context("with auth disabled", func() {
		BeforeEach(func() {
			svc = NewService(name, true)
		})

		It("should inject an authenticated user", func() {
			RunHandler(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
				Ω(req.Header.Get(AuthorizationHeader)).ShouldNot(BeEmpty())
			})
		})

	})

	BeforeEach(func() {
		svc = NewService(name, false)
	})

	It("should not have a user injected", func() {
		RunHandler(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			Ω(req.Header.Get(AuthorizationHeader)).Should(BeEmpty())
		})
	})

	It("should register request ID middleware", func() {
		RunHandler(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			Ω(middleware.ContextRequestID(ctx)).ShouldNot(BeEmpty())
		})
	})

	It("should log requests", func() {
		var reqid string
		RunHandler(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			reqid = middleware.ContextRequestID(ctx)
		})
		Ω(reqid).ShouldNot(BeEmpty())
		Ω(logBuf.String()).Should(ContainSubstring(fmt.Sprintf("started req_id=%s", reqid)))
		Ω(logBuf.String()).Should(ContainSubstring(fmt.Sprintf("completed req_id=%s", reqid)))
	})

	It("should register a metric registry", func() {
		RunHandler(func(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
			Ω(ContextMetrics(ctx)).ShouldNot(BeNil())
		})
	})

	It("should recover from panics", func() {
		Ω(func() {
			handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
				panic("o no")
				return nil
			}
			ctrl.MuxHandler("test", handler, nil)(rw, req, url.Values{})
		}).ShouldNot(Panic())
	})

	It("should log uncaught errors", func() {
		errstr := RandStringRunes(8)
		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			return errors.New(errstr)
		}
		ctrl.MuxHandler("test", handler, nil)(rw, req, url.Values{})
		Ω(logBuf.String()).Should(ContainSubstring(fmt.Sprintf("err=%s", errstr)))
		Ω(rw.Code).Should(Equal(500))
	})
})
