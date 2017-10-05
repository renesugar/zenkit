package admin

import "github.com/zenoss/zenkit/healthcheck"

var updater healthcheck.Updater

func ResetRegistry() {
	healthcheck.DefaultRegistry = healthcheck.NewRegistry()
	updater = healthcheck.NewStatusUpdater()
	healthcheck.Register("manual_http_status", updater)
}

func init() {
	ResetRegistry()
}
