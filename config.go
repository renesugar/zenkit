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

	SignKeyFileConfig = "sign.key_file"
	SignMethodConfig  = "sign.method"
	SignExpiryConfig  = "sign.expiry"

	HTTPPortConfig  = "http.port"
	AdminPortConfig = "admin.port"

	GCProjectIDConfig         = "gcloud.project_id"
	GCNoAuthConfig            = "gcloud.no_auth"
	GCEmulatorBigtableConfig  = "gcloud.emulator.bigtable"
	GCEmulatorDatastoreConfig = "gcloud.emulator.datastore"
	GCEmulatorPubsubConfig    = "gcloud.emulator.pubsub"
)

func AddStandardServerOptions(cmd *cobra.Command, port, adminPort int) {
	AddHTTPOptions(cmd, port)
	AddAdminOptions(cmd, adminPort)
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
	cmd.PersistentFlags().String("auth-key-file", "/run/secrets/auth_key", "File containing authentication verification key")
	viper.BindPFlag(AuthKeyFileConfig, cmd.PersistentFlags().Lookup("auth-key-file"))
	viper.SetDefault(AuthKeyFileConfig, "/run/secrets/auth_key")

	cmd.PersistentFlags().Bool("auth-disabled", false, "Run with middleware that injects a default admin identity for unauthenticated requests")
	viper.BindPFlag(AuthDisabledConfig, cmd.PersistentFlags().Lookup("auth-disabled"))
	viper.SetDefault(AuthDisabledConfig, false)
}

func AddSignConfigOptions(cmd *cobra.Command) {
	cmd.PersistentFlags().String("sign-key-file", "/run/secrets/sign_key", "File containing key to sign jwt with")
	viper.BindPFlag(SignKeyFileConfig, cmd.PersistentFlags().Lookup("sign-key-file"))
	viper.SetDefault(SignKeyFileConfig, "/run/secrets/sign_key")

	cmd.PersistentFlags().String("sign-method", "HS256", "Method used to sign jwt")
	viper.BindPFlag(SignMethodConfig, cmd.PersistentFlags().Lookup("sign-method"))
	viper.SetDefault(SignMethodConfig, "HS256")

	cmd.PersistentFlags().Int("sign-expiry", 60, "Duration, in seconds, a signed jwt is valid for")
	viper.BindPFlag(SignExpiryConfig, cmd.PersistentFlags().Lookup("sign-expiry"))
	viper.SetDefault(SignExpiryConfig, 60)
}

func AddHTTPOptions(cmd *cobra.Command, port int) {
	cmd.PersistentFlags().IntP("http-port", "p", port, "Port to which the server should bind")
	viper.BindPFlag(HTTPPortConfig, cmd.PersistentFlags().Lookup("http-port"))
	viper.SetDefault(HTTPPortConfig, fmt.Sprintf("%d", port))
}

func AddAdminOptions(cmd *cobra.Command, adminPort int) {
	cmd.PersistentFlags().Int("admin-port", adminPort, "Port to which the admin server should bind")
	viper.BindPFlag(AdminPortConfig, cmd.PersistentFlags().Lookup("admin-port"))
	viper.SetDefault(AdminPortConfig, fmt.Sprintf("%d", adminPort))
}

func AddGCloudOptions(cmd *cobra.Command) {
	cmd.PersistentFlags().String("gcloud-project-id", "", "Google Cloud project/dataset id")
	viper.BindPFlag(GCProjectIDConfig, cmd.PersistentFlags().Lookup("gcloud-project-id"))
	viper.SetDefault(GCProjectIDConfig, "")

	cmd.PersistentFlags().Bool("gcloud-no-auth", false, "Disable Google Cloud auth")
	viper.BindPFlag(GCNoAuthConfig, cmd.PersistentFlags().Lookup("gcloud-no-auth"))
	viper.SetDefault(GCNoAuthConfig, false)
}

func AddGCloudEmulatorBigtableOptions(cmd *cobra.Command) {
	cmd.PersistentFlags().String("gcloud-emulator-bigtable", "", "Host:port of the gcloud bigtable emulator")
	viper.BindPFlag(GCEmulatorBigtableConfig, cmd.PersistentFlags().Lookup("gcloud-emulator-bigtable"))
	viper.SetDefault(GCEmulatorBigtableConfig, "")
}

func AddGCloudEmulatorDatastoreOptions(cmd *cobra.Command) {
	cmd.PersistentFlags().String("gcloud-emulator-datastore", "", "Host:port of the gcloud datastore emulator")
	viper.BindPFlag(GCEmulatorDatastoreConfig, cmd.PersistentFlags().Lookup("gcloud-emulator-datastore"))
	viper.SetDefault(GCEmulatorDatastoreConfig, "")
}

func AddGCloudEmulatorPubsubOptions(cmd *cobra.Command) {
	cmd.PersistentFlags().String("gcloud-emulator-pubsub", "", "Host:port of the gcloud pubsub emulator")
	viper.BindPFlag(GCEmulatorPubsubConfig, cmd.PersistentFlags().Lookup("gcloud-emulator-pubsub"))
	viper.SetDefault(GCEmulatorPubsubConfig, "")
}
