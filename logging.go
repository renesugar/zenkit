package zenkit

import (
	"context"

	"github.com/Sirupsen/logrus"
	"github.com/goadesign/goa"
	goalogrus "github.com/goadesign/goa/logging/logrus"
)

var noop = func() {}

func ServiceLogger() goa.LogAdapter {
	return goalogrus.New(logrus.New())
}

func ContextLogger(ctx context.Context) *logrus.Entry {
	return goalogrus.Entry(ctx)
}

func SetVerbosity(svc *goa.Service, verbosity int) {
	logger := ContextLogger(svc.Context).Logger
	logger.Level = logrus.WarnLevel + logrus.Level(verbosity)
	logger.WithField("level", logger.Level).Info("Log level changed")
}

func LogEntryAndExit(ctx context.Context) func() {
	logger := ContextLogger(ctx)
	if logger == nil {
		return noop
	}
	fn := funcName(2)
	logger.Debugf("ENTER %s()", fn)
	exit := func() {
		logger.Debugf("EXIT %s()", fn)
	}
	return exit
}
