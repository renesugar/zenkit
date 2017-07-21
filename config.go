package zenkit

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	LogLevelConfig = "log.level"

	TracingEnabledConfig    = "tracing.enabled"
	TracingDaemonConfig     = "tracing.daemon"
	TracingSampleRateConfig = "tracing.sample_rate"

	AuthDisabledConfig = "auth.disabled"
	AuthKeyFileConfig  = "auth.key_file"

	HTTPPortConfig = "http.port"
)

func AddStandardServerOptions(cmd *cobra.Command, port int) {
	AddHTTPOptions(cmd, port)
	AddAuthConfigOptions(cmd)
	AddTracingConfigOptions(cmd)
}

func AddLoggingConfigOptions(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("log-level", "v", "info", "Log level")
	viper.BindPFlag(LogLevelConfig, cmd.PersistentFlags().Lookup("log-level"))
	viper.SetDefault(LogLevelConfig, "info")
}

func AddTracingConfigOptions(cmd *cobra.Command) {
	cmd.PersistentFlags().Bool("tracing-enabled", false, "Whether to send trace info to AWS X-Ray")
	viper.BindPFlag(TracingEnabledConfig, cmd.PersistentFlags().Lookup("tracing-enabled"))
	viper.SetDefault(TracingEnabledConfig, false)

	cmd.PersistentFlags().String("tracing-daemon", "", "Address of the AWS X-Ray daemon")
	viper.BindPFlag(TracingDaemonConfig, cmd.PersistentFlags().Lookup("tracing-daemon"))
	viper.SetDefault(TracingDaemonConfig, "")

	cmd.PersistentFlags().Int("tracing-sample-rate", 100, "Rate at which tracing should sample requests")
	viper.BindPFlag(TracingSampleRateConfig, cmd.PersistentFlags().Lookup("tracing-sample-rate"))
	viper.SetDefault(TracingSampleRateConfig, 100)
}

func AddAuthConfigOptions(cmd *cobra.Command) {
	cmd.PersistentFlags().String("auth-key-file", "", "File containing authentication verification key")
	viper.BindPFlag(AuthKeyFileConfig, cmd.PersistentFlags().Lookup("auth-key-file"))
	viper.SetDefault(AuthKeyFileConfig, "")

	cmd.PersistentFlags().Bool("auth-disabled", false, "Run with middleware that injects a default admin identity for unauthenticated requests")
	viper.BindPFlag(AuthDisabledConfig, cmd.PersistentFlags().Lookup("auth-disabled"))
	viper.SetDefault(AuthDisabledConfig, false)
}

func AddHTTPOptions(cmd *cobra.Command, port int) {
	cmd.PersistentFlags().IntP("http-port", "p", port, "Port to which the server should bind")
	viper.BindPFlag(HTTPPortConfig, cmd.PersistentFlags().Lookup("http-port"))
	viper.SetDefault(HTTPPortConfig, fmt.Sprintf("%d", port))
}
