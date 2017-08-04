package zenkit_test

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/goatest"
	. "github.com/zenoss/zenkit"
	"github.com/zenoss/zenkit/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("Tracing", func() {

	var (
		address    string
		listener   *net.UDPConn
		buffer     *Buffer
		svc        *goa.Service
		req        *http.Request
		rw         *httptest.ResponseRecorder
		ctrl       *goa.Controller
		logBuf     *Buffer
		logger     *log.Logger
		resp       interface{}
		respSetter goatest.ResponseSetterFunc = func(r interface{}) { resp = r }
		encoder                               = func(io.Writer) goa.Encoder { return respSetter }
	)

	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	listener, _ = net.ListenUDP("udp", addr)
	address = listener.LocalAddr().String()
	listener.SetReadDeadline(time.Now().Add(time.Second))
	buffer = BufferReader(listener)

	BeforeEach(func() {
		svc = NewService(test.RandString(8), false)

		logBuf = NewBuffer()
		logger = log.New(logBuf, "", log.Ltime)

		req, _ = http.NewRequest("", "http://example.com/", nil)
		rw = httptest.NewRecorder()

		svc.WithLogger(goa.NewLogger(logger))
		svc.Encoder = goa.NewHTTPEncoder()
		svc.Encoder.Register(encoder, "*/*")
		ctrl = svc.NewController("test")

	})

	Describe("middleware", func() {
		It("should be able to be used by a service", func() {
			err := UseXRayMiddleware(svc, address, 100)
			Ω(err).ShouldNot(HaveOccurred())

			handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
				// Nothing to see here
				return nil
			}
			ctrl.MuxHandler("test", handler, nil)(rw, req, url.Values{})
			Eventually(buffer).Should(Say(`{"format": "json", "version": 1}`)) // This is the header of a trace packet
		})

		It("should fail if there's a bad address", func() {
			err := UseXRayMiddleware(svc, "notanaddress", 100)
			Ω(err).Should(HaveOccurred())
		})

		It("should fail if there is no address", func() {
			err := UseXRayMiddleware(svc, "", 100)
			Ω(err).Should(HaveOccurred())
		})
	})

})
