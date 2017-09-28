package metrics_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	. "github.com/zenoss/zenkit/metrics"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TimedFunc(ctx context.Context) {
	defer MeasureTime(ctx)()
	time.Sleep(1 * time.Millisecond)
}

var _ = Describe("Metrics", func() {

	It("should be able to be registered on and retrieved from the context", func() {
		ctx := context.Background()
		reg := metrics.NewRegistry()
		newctx := WithMetrics(ctx, reg)
		newreg := ContextMetrics(newctx)
		Ω(newreg).Should(BeIdenticalTo(reg))
	})

	It("should not fail if no metric registry exists", func() {
		ctx := context.Background()
		reg := ContextMetrics(ctx)
		Ω(reg).Should(BeNil())
	})

	It("should be able to be registered by middleware", func() {
		var (
			m      metrics.Registry
			resp   = httptest.NewRecorder()
			req, _ = http.NewRequest("", "http://example.com", nil)
		)

		handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			m = ContextMetrics(ctx)
			return nil
		}

		mw := MetricsMiddleware()
		wrapped := mw(handler)
		err := wrapped(context.Background(), resp, req)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(m).ShouldNot(BeNil())
	})

	It("should be able to measure the time a function takes", func() {
		reg := metrics.NewRegistry()
		ctx := WithMetrics(context.Background(), reg)
		TimedFunc(ctx)
		metric := reg.Get("func.TimedFunc.time")
		Ω(metric).ShouldNot(BeNil())
	})

	It("should not panic function metrics if no registry is defined", func() {
		TimedFunc(context.Background())
	})

	It("should be able to increment a counter", func() {
		reg := metrics.NewRegistry()
		ctx := WithMetrics(context.Background(), reg)
		IncrementCounter(ctx, "my.test.counter", 1)
		metric := reg.Get("my.test.counter")
		Ω(metric).ShouldNot(BeNil())
	})

	It("should not panic counter increment if no registry is defined", func() {
		IncrementCounter(context.Background(), "my.test.counter", 1)
	})
})
