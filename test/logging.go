package test

import (
	"github.com/Sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	// logger.Out = ginkgo.GinkgoWriter
}

func TestLogger() *logrus.Logger {
	return logger
}
