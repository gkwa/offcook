package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/gkwa/offcook/core"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/gkwa/offcook/internal/logger"
)

var (
	cfgFile   string
	verbose   bool
	logFormat string
	cliLogger logr.Logger
)

var rootCmd = &cobra.Command{
	Use:   "offcook",
	Short: "offcook is a tool for managing environment variables",
	Long:  `offcook generates scripts for various shells and tools to set environment variables.`,
	Run: func(cmd *cobra.Command, args []string) {
		vars, _ := cmd.Flags().GetStringSlice("var")

		logger := LoggerFrom(cmd.Context())
		executor := core.NewTemplateExecutor(logger)

		err := executor.ExecuteToWriter(vars, os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringSlice("var", []string{}, "List of environment variables")
	if err := rootCmd.MarkFlagRequired("var"); err != nil {
		fmt.Printf("Error marking 'var' flag as required: %v\n", err)
		os.Exit(1)
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.offcook.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose mode")
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "", "json or text (default is text)")

	if err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		fmt.Printf("Error binding verbose flag: %v\n", err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("log-format", rootCmd.PersistentFlags().Lookup("log-format")); err != nil {
		fmt.Printf("Error binding log-format flag: %v\n", err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".offcook")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	logFormat = viper.GetString("log-format")
	verbose = viper.GetBool("verbose")
}

func LoggerFrom(ctx context.Context, keysAndValues ...interface{}) logr.Logger {
	if cliLogger.IsZero() {
		cliLogger = logger.NewConsoleLogger(verbose, logFormat == "json")
	}
	newLogger := cliLogger
	if ctx != nil {
		if l, err := logr.FromContext(ctx); err == nil {
			newLogger = l
		}
	}
	return newLogger.WithValues(keysAndValues...)
}
