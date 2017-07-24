package zenkit_test

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	. "github.com/zenoss/zenkit"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	var (
		prefix string
		cmd    *cobra.Command
		port   int
	)

	setenv := func(s, v string) {
		varname := fmt.Sprintf("%s_%s", strings.ToUpper(prefix), s)
		os.Setenv(varname, v)
	}

	BeforeEach(func() {
		prefix = RandStringRunes(8)
		cmd = &cobra.Command{Use: "c", Run: func(*cobra.Command, []string) {}}
		viper.Reset()
		viper.AutomaticEnv()
		viper.SetEnvPrefix(prefix)
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	})

	TestLoggingFlags := func() {
		It("should have a default log level of info", func() {
			Ω(viper.Get(LogLevelConfig)).Should(Equal("info"))
		})

		It("should accept configuration of log level via env var", func() {
			viper.AutomaticEnv()
			setenv("LOG_LEVEL", "debug")
			Ω(viper.Get(LogLevelConfig)).Should(Equal("debug"))
		})
		It("should accept configuration of log level via command line", func() {
			err := cmd.ParseFlags([]string{"--log-level", "error"})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(viper.Get(LogLevelConfig)).Should(Equal("error"))
		})
	}

	TestHTTPFlags := func() {
		It("should have a default port based on what's passed in", func() {
			Ω(viper.GetInt(HTTPPortConfig)).Should(BeNumerically("==", port))
		})

		It("should allow setting the port via env var", func() {
			port2 := port - 1
			setenv("HTTP_PORT", fmt.Sprintf("%d", port2))
			Ω(viper.GetInt(HTTPPortConfig)).Should(BeNumerically("==", port2))
		})

		It("should allow setting the port via command line", func() {
			port2 := port - 1
			err := cmd.ParseFlags([]string{"--http-port", fmt.Sprintf("%d", port2)})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(viper.GetInt(HTTPPortConfig)).Should(BeNumerically("==", port2))
		})
	}

	TestTracingFlags := func() {
		It("should have tracing disabled by default", func() {
			Ω(viper.GetBool(TracingEnabledConfig)).Should(BeFalse())
		})

		It("should have no daemon set by default", func() {
			Ω(viper.GetString(TracingDaemonConfig)).Should(BeEmpty())
		})

		It("should have a sample rate of 100 by default", func() {
			Ω(viper.GetInt(TracingSampleRateConfig)).Should(BeNumerically("==", 100))
		})

		It("should allow enabling tracing via env var", func() {
			setenv("TRACING_ENABLED", "1")
			Ω(viper.GetBool(TracingEnabledConfig)).Should(BeTrue())
		})

		It("should allow setting the tracing daemon via env var", func() {
			daemon := RandStringRunes(10)
			setenv("TRACING_DAEMON", daemon)
			Ω(viper.GetString(TracingDaemonConfig)).Should(Equal(daemon))
		})

		It("should allow setting the tracing sample rate via env var", func() {
			n := rand.Intn(1000)
			setenv("TRACING_SAMPLE_RATE", fmt.Sprintf("%d", n))
			Ω(viper.GetInt(TracingSampleRateConfig)).Should(BeNumerically("==", n))
		})

		It("should allow enabling tracing via command line", func() {
			err := cmd.ParseFlags([]string{"--tracing-enabled"})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(viper.GetBool(TracingEnabledConfig)).Should(BeTrue())
		})

		It("should allow setting the tracing daemon via command line", func() {
			daemon := RandStringRunes(10)
			err := cmd.ParseFlags([]string{"--tracing-daemon", daemon})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(viper.GetString(TracingDaemonConfig)).Should(Equal(daemon))
		})

		It("should allow setting the tracing sample rate via command line", func() {
			n := rand.Intn(1000)
			err := cmd.ParseFlags([]string{"--tracing-sample-rate", fmt.Sprintf("%d", n)})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(viper.GetInt(TracingSampleRateConfig)).Should(BeNumerically("==", n))
		})

	}

	TestAuthFlags := func() {
		It("should not have a default auth key file", func() {
			Ω(viper.Get(AuthKeyFileConfig)).Should(BeEmpty())
		})

		It("should have auth enabled by default", func() {
			Ω(viper.Get(AuthDisabledConfig)).Should(BeFalse())
		})

		It("should allow setting the key file via env var", func() {
			keyfile := RandStringRunes(10)
			setenv("AUTH_KEY_FILE", keyfile)
			Ω(viper.GetString(AuthKeyFileConfig)).Should(Equal(keyfile))
		})

		It("should allow disabling auth via env var", func() {
			setenv("AUTH_DISABLED", "1")
			Ω(viper.GetBool(AuthDisabledConfig)).Should(BeTrue())
		})

		It("should allow setting the key file via command line", func() {
			keyfile := RandStringRunes(10)
			err := cmd.ParseFlags([]string{"--auth-key-file", keyfile})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(viper.Get(AuthKeyFileConfig)).Should(Equal(keyfile))
		})

		It("should allow disabling auth via command line", func() {
			err := cmd.ParseFlags([]string{"--auth-disabled"})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(viper.GetBool(AuthDisabledConfig)).Should(BeTrue())
		})
	}

	Context("with tracing flags", func() {

		BeforeEach(func() {
			AddTracingConfigOptions(cmd)
		})

		TestTracingFlags()
	})

	Context("with auth flags", func() {
		BeforeEach(func() {
			AddAuthConfigOptions(cmd)
		})

		TestAuthFlags()

	})

	Context("with HTTP flags", func() {

		BeforeEach(func() {
			port = rand.Intn(65535)
			AddHTTPOptions(cmd, port)
		})

		TestHTTPFlags()

	})

	Context("with logging flags", func() {

		BeforeEach(func() {
			AddLoggingConfigOptions(cmd)
		})

		TestLoggingFlags()

	})

	Context("with standard server flags", func() {

		BeforeEach(func() {
			port = rand.Intn(65535)
			AddStandardServerOptions(cmd, port)
		})

		TestHTTPFlags()
		TestAuthFlags()
		TestTracingFlags()
	})

})
