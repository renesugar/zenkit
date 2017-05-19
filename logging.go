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

func LogEntryAndExit(ctx context.Context) func() {
	logger := ContextLogger(ctx)
	fn := funcName(3)
	logger.Debugf("ENTER %s()", fn)
	exit := func() {
		log.Debugf("EXIT %s()", fn)
	}
	return exit
}
