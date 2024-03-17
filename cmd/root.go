package cmd

import (
	"os"
	"strings"

	"github.com/fortix/galera-ips/internal/monitor"
	"github.com/fortix/galera-ips/internal/server"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const CONFIG_FILE_NAME = "galera-ips"
const CONFIG_FILE_TYPE = "yaml"
const CONFIG_ENV_PREFIX = "GALERA_IPS"

var rootCmd = &cobra.Command{
	Use:   "galera-ips",
	Short: "Expose Galera Cluster IPs via REST API",
	Long:  `This command periodically queries ProxySQL for the list of IPs belonging to the read and write groups.`,
	Args:  cobra.NoArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("proxysql.host", cmd.Flags().Lookup("host"))
		viper.BindEnv("proxysql.host", CONFIG_ENV_PREFIX+"_HOST")
		viper.SetDefault("proxysql.host", "127.0.0.1")

		viper.BindPFlag("proxysql.port", cmd.Flags().Lookup("port"))
		viper.BindEnv("proxysql.port", CONFIG_ENV_PREFIX+"_PORT")
		viper.SetDefault("proxysql.port", 6032)

		viper.BindPFlag("proxysql.user", cmd.Flags().Lookup("user"))
		viper.BindEnv("proxysql.user", CONFIG_ENV_PREFIX+"_USER")
		viper.SetDefault("proxysql.user", "")

		viper.BindPFlag("proxysql.password", cmd.Flags().Lookup("password"))
		viper.BindEnv("proxysql.password", CONFIG_ENV_PREFIX+"_PASSWORD")
		viper.SetDefault("proxysql.password", "")

		viper.BindPFlag("listen", cmd.Flags().Lookup("listen"))
		viper.BindEnv("listen", CONFIG_ENV_PREFIX+"_LISTEN")
		viper.SetDefault("listen", "0.0.0.0:6031")
	},
	Run: serveCmd,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("config", "c", "", "Config file (default is "+CONFIG_FILE_NAME+"."+CONFIG_FILE_TYPE+" in the current directory or $HOME/).\nOverrides the "+CONFIG_ENV_PREFIX+"_CONFIG environment variable if set.")
	rootCmd.PersistentFlags().StringP("log-level", "", "info", "Log level (debug, info, warn, error, fatal, panic).\nOverrides the "+CONFIG_ENV_PREFIX+"_LOGLEVEL environment variable if set.")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
		viper.BindEnv("config", CONFIG_ENV_PREFIX+"_CONFIG")
		viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))
		viper.BindEnv("log.level", CONFIG_ENV_PREFIX+"_LOGLEVEL")

		// If config file given then use it
		cfgFile := viper.GetString("config")
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
			if err := viper.ReadInConfig(); err != nil {
				log.Fatal().Msgf("missing config file: %s", viper.ConfigFileUsed())
			}
		}

		switch viper.GetString("log.level") {
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "info":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "warn":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case "error":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		default:
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		}
	}

	rootCmd.Flags().StringP("host", "", "127.0.0.1", "The address of the ProxySQL server to connect to.\nOverrides the "+CONFIG_ENV_PREFIX+"_HOST environment variable if set.")
	rootCmd.Flags().IntP("port", "", 6032, "The port of the ProxySQL server to connect to.\nOverrides the "+CONFIG_ENV_PREFIX+"_PORT environment variable if set.")
	rootCmd.Flags().StringP("username", "", "", "The username to use when connecting to the ProxySQL server.\nOverrides the "+CONFIG_ENV_PREFIX+"_USER environment variable if set.")
	rootCmd.Flags().StringP("password", "", "", "The password to use when connecting to the ProxySQL server.\nOverrides the "+CONFIG_ENV_PREFIX+"_PASSWORD environment variable if set.")
	rootCmd.Flags().StringP("listen", "l", "0.0.0.0:6031", "The address and port to listen on.\nOverrides the "+CONFIG_ENV_PREFIX+"_LISTEN environment variable if set.")
}

func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Set search paths for config file
	viper.AddConfigPath(".")
	viper.AddConfigPath(home)
	viper.AddConfigPath("/etc")
	viper.SetConfigName(CONFIG_FILE_NAME) // Name of config file without extension
	viper.SetConfigType(CONFIG_FILE_TYPE) // Type of config file
	viper.SetEnvPrefix(CONFIG_ENV_PREFIX)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	viper.ReadInConfig()
}

func serveCmd(cmd *cobra.Command, args []string) {

	// Check the ProxySQL user and password are set
	if viper.GetString("proxysql.username") == "" || viper.GetString("proxysql.password") == "" {
		log.Fatal().Msgf("You must provide a ProxySQL user and password")
	}

	go monitor.MonitorProxySQL()

	server.Run()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
