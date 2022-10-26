package cmd

import (
	"context"
	"fmt"
	"syscall"

	"github.com/meimeitou/makabaka/config"
	"github.com/meimeitou/makabaka/server"
	"github.com/oklog/run"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type serveOptions struct {
	// Config file path
	config string

	// Flags
	webHTTPAddr  string
	webHTTPSAddr string
	webPrefix    string
}

func commandServe() *cobra.Command {
	options := serveOptions{}

	cmd := &cobra.Command{
		Use:     "serve [flags] [config file]",
		Short:   "Run makabaka",
		Example: "makabaka serve config.yaml",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
			options.config = args[0]
			return runServe(options, logrus.StandardLogger())
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&options.webHTTPAddr, "web-http-addr", "", "Web HTTP address")
	flags.StringVar(&options.webHTTPSAddr, "web-https-addr", "", "Web HTTPS address")
	flags.StringVar(&options.webPrefix, "web-prefix", "", "Web HTTPS address")

	return cmd
}

func runServe(options serveOptions, logger *logrus.Logger) error {
	// init config
	configFile := options.config
	c, err := config.ReadConfigFile(configFile)
	if err != nil {
		return err
	}
	level, err := logrus.ParseLevel(c.Log)
	if err != nil {
		return err
	}
	logger.SetLevel(level)
	applyConfigOverrides(options, c)
	logrus.Info(c)
	// init storage
	err = config.InitDBSet(c, logger)
	if err != nil {
		return err
	}
	// run all
	config.G.Add(run.SignalHandler(context.Background(), syscall.SIGTERM, syscall.SIGINT))
	sv := server.NewServer(logger, c.Server.HTTP, c.Server.Prefix)
	sv.Run(&config.G)
	if err := config.G.Run(); err != nil {
		if _, ok := err.(run.SignalError); !ok {
			return fmt.Errorf("run groups: %w", err)
		}
		logrus.Infof("%v, shutdown now", err)
	}
	return nil
}

func applyConfigOverrides(options serveOptions, cfg *config.Config) {
	if options.webHTTPAddr != "" {
		cfg.Server.HTTP = options.webHTTPAddr
	}

	if options.webHTTPSAddr != "" {
		cfg.Server.HTTPS = options.webHTTPSAddr
	}

	if options.webPrefix != "" {
		cfg.Server.Prefix = options.webPrefix
	}

}

func init() {
	rootCmd.AddCommand(commandServe())
}
