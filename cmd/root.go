package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/webitel/wlog"

	"github.com/kirychukyurii/fd-import/config"
)

var (

	// version is the app's semantic version.
	version = "0.0.0"

	// commit is the git commit used to build the App.
	commit     = "hash"
	commitDate = "date"
)

func Execute() {

	// os.Interrupt to gracefully shutdown on Ctrl+C which is SIGINT
	// syscall.SIGTERM is the usual signal for termination and the default one (it can be modified)
	// for docker containers, which is also used by kubernetes.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := command().ExecuteContext(ctx); err != nil {
		os.Exit(-1)
	}
}

func command() *cobra.Command {
	log := wlog.NewLogger(&wlog.LoggerConfiguration{
		EnableConsole: true,
		ConsoleLevel:  wlog.LevelDebug,
	})

	cfg := config.New()
	c := &cobra.Command{
		Use:          "fd-import",
		Short:        "FreshDesk Import - easy import .json files, exported from API",
		SilenceUsage: true,
		Version:      fmt.Sprintf("%s, commit %s, date %s", version, commit, commitDate),
		Args: func(cmd *cobra.Command, args []string) error {

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// apply CLI args to config
			if err := cmd.ParseFlags(os.Args[1:]); err != nil {
				return fmt.Errorf("parsing flags: %w", err)
			}

			return nil
		},
	}

	flagSet(c.PersistentFlags(), cfg)
	c.AddCommand(importCommand(cfg, log), migrateCommand(cfg, log), apiCommand(cfg, log))

	return c
}

// flagSet sets up the command line flags using the given FlagSet and Config objects.
// The function binds the flags to the corresponding fields in the Config object.
// The flag names, shorthand flags, default values, and usage descriptions are specified
// for each flag.
func flagSet(fs *pflag.FlagSet, cfg *config.Config) {
	fs.StringVarP(&cfg.LogLevel, "log-level", "l", "debug", "log level")
	fs.StringVar(&cfg.LogFile, "log-file", "./fd-import.log", "log file")
	fs.StringVarP(&cfg.DSN, "dsn", "d", "", "database connection string")
}
