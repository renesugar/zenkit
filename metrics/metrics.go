package metrics

import (
	"context"
	"fmt"

	metrics "github.com/rcrowley/go-metrics"
	"github.com/zenoss/zenkit/funcname"
)

func MeasureTime(ctx context.Context) func() {
	begin := TimeFunc()
	fn := funcname.FuncName(2)
	registry := ContextMetrics(ctx)
	exit := func() {
		if registry == nil {
			return
		}
		t := metrics.GetOrRegisterTimer(fmt.Sprintf("func.%s.time", fn), registry)
		t.UpdateSince(begin)
	}
	return exit
}

func IncrementCounter(ctx context.Context, name string, inc int64) {
	registry := ContextMetrics(ctx)
	if registry == nil {
		return
	}
	ctr := metrics.GetOrRegisterCounter(name, registry)
	ctr.Inc(inc)
}
