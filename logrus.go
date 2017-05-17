package zenkit

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/goadesign/goa"
	goalogrus "github.com/goadesign/goa/logging/logrus"
)

func ServiceLogger() goa.LogAdapter {
	return goalogrus.New(logrus.New())
}

func ContextLogger(ctx context.Context) *logrus.Entry {
	return goalogrus.Entry(ctx)
}

func SetVerbosity(svc *goa.Service, verbosity int) {
	logger := ContextLogger(svc.Context).Logger
	logger.Level = logrus.WarnLevel + logrus.Level(verbosity)
}
