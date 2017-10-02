package logging

import (
	"context"

	"github.com/goadesign/goa"
	goalogrus "github.com/goadesign/goa/logging/logrus"
	"github.com/sirupsen/logrus"
	"github.com/zenoss/zenkit/funcname"
)

var noop = func() {}

type ErrorLogger interface {
	LogError(msg string, keys ...interface{})
}

func ServiceLogger() goa.LogAdapter {
	return goalogrus.New(logrus.New())
}

func ContextLogger(ctx context.Context) *logrus.Entry {
	return goalogrus.Entry(ctx)
}

func SetLogLevel(svc *goa.Service, level string) {
	logger := ContextLogger(svc.Context).Logger
	oldlevel := logger.Level
	newlevel, err := logrus.ParseLevel(level)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"level":    oldlevel,
			"badlevel": level,
		}).Error("Unable to parse log level. Not changing.")
		return
	}
	if oldlevel == newlevel {
		logger.WithFields(logrus.Fields{
			"level": oldlevel,
		}).Debug("Requested log level is already active. Ignoring.")
		return
	}
	logger.Level = newlevel
	logger.WithFields(logrus.Fields{
		"oldlevel": oldlevel,
		"newlevel": newlevel,
	}).Info("Log level changed")
}

func LogEntryAndExit(ctx context.Context) func() {
	logger := ContextLogger(ctx)
	if logger == nil {
		return noop
	}
	fn := funcname.FuncName(2)
	logger.Debugf("ENTER %s()", fn)
	exit := func() {
		logger.Debugf("EXIT %s()", fn)
	}
	return exit
}
