package checks_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	. "github.com/zenoss/zenkit/healthcheck/checks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Checks", func() {

	Context("with a FileChecker", func() {

		var (
			d   string
			err error
		)

		BeforeEach(func() {
			d, err = ioutil.TempDir("", "test-file-checker")
			Ω(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			if d != "" {
				err = os.RemoveAll(d)
				Ω(err).ShouldNot(HaveOccurred())
			}
		})

		It("should return nil if the file does not exist", func() {
			f := filepath.Join(d, "test-file")

			c := FileChecker(f)
			Ω(c.Check()).Should(BeNil())

			By("creating the file")

			_, err := os.Create(f)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(c.Check()).ShouldNot(BeNil())
		})
	})

	Context("with an HTTPChecker", func() {
		var (
			server *ghttp.Server
		)

		BeforeEach(func() {
			server = ghttp.NewServer()
		})

		AfterEach(func() {
			server.Close()
		})

		It("should return an error if the url is invalid", func() {
			c := HTTPChecker("#%", 200, time.Second, http.Header{})
			Ω(c.Check()).Should(HaveOccurred())
		})

		It("should return an error if it cannot connect to the service", func() {
			c := HTTPChecker("abc123", 200, time.Second, http.Header{})
			Ω(c.Check()).Should(HaveOccurred())
		})

		It("should return an error if the status code does not match", func() {
			statusCode := http.StatusServiceUnavailable
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("HEAD", "/"),
				ghttp.RespondWith(statusCode, nil),
			))
			c := HTTPChecker(server.URL(), 200, time.Second, http.Header{})
			Ω(c.Check()).Should(HaveOccurred())
		})

		It("should return nil if the status code does match", func() {
			h := http.Header{}
			h.Set("test", "val")
			statusCode := http.StatusOK
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("HEAD", "/"),
				ghttp.VerifyHeader(h),
				ghttp.RespondWith(statusCode, nil),
			))
			c := HTTPChecker(server.URL(), 200, time.Second, h)
			Ω(c.Check()).ShouldNot(HaveOccurred())
		})
	})

	Context("with a TCPChecker", func() {
		var (
			server *ghttp.Server
		)

		BeforeEach(func() {
			server = ghttp.NewServer()
		})

		AfterEach(func() {
			server.Close()
		})

		It("should return an error if it cannot connect to the service", func() {
			c := TCPChecker("abc123", time.Second)
			Ω(c.Check()).Should(HaveOccurred())
		})

		It("should return nil if it connects to the service", func() {
			c := TCPChecker(server.Addr(), time.Second)
			Ω(c.Check()).ShouldNot(HaveOccurred())
		})
	})
})
