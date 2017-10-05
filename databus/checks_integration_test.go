//+build integration

package databus_test

import (
	"time"

	. "github.com/zenoss/zenkit/databus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("ChecksIntegration", func() {

	Context("with a KafkaChecker", func() {

		var (
			addr string
			err  error
		)

		BeforeEach(func() {
			addr, err = harness.Resolve("kafka", 9092)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should return an error if it cannot connect to the kafka client", func() {
			c := KafkaChecker("dummy")
			Ω(c.Check()).Should(HaveOccurred())
		})

		It("should return a success if it can connect to the kafka client", func() {
			c := KafkaChecker(addr)
			Ω(c.Check()).ShouldNot(HaveOccurred())
		})
	})

	Context("with a SchemaRegistryChecker", func() {
		var (
			addr string
			err  error
		)

		BeforeEach(func() {
			addr, err = harness.Resolve("kafka-schema-registry", 8081)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should return an error if the url is invalid", func() {
			c := SchemaRegistryChecker("#%", time.Second)
			Ω(c.Check()).Should(HaveOccurred())
		})

		It("should return an error if it cannot connect to the sr-client", func() {
			c := SchemaRegistryChecker("dummy", time.Second)
			Ω(c.Check()).Should(HaveOccurred())
		})

		It("should return an error if the response body is invalid", func() {
			server := ghttp.NewServer()
			server.AllowUnhandledRequests = true
			defer server.Close()
			c := SchemaRegistryChecker(server.Addr(), time.Second)
			Ω(c.Check()).Should(HaveOccurred())
		})

		It("should return a success if the response body is valid", func() {
			c := SchemaRegistryChecker(addr, time.Second)
			Ω(c.Check()).ShouldNot(HaveOccurred())
		})
	})
})
