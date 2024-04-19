package cmd

import (
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/webitel/wlog"

	"github.com/kirychukyurii/fd-import/config"
	"github.com/kirychukyurii/fd-import/pkg/db"
)

func migrateCommand(cfg *config.Config, log *wlog.Logger) *cobra.Command {
	var directory string

	c := &cobra.Command{
		Use:          "migrate",
		Short:        "Apply database schema migrations",
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

			pool, err := db.New(cmd.Context(), log, cfg.DSN)
			if err != nil {
				return err
			}

			goose.SetTableName(fmt.Sprintf("fresh.%s", goose.DefaultTablename))
			if err := goose.UpContext(cmd.Context(), pool.STDLib(), directory); err != nil {
				return err
			}

			return nil
		},
	}

	directory = migrateFlagSet(c.PersistentFlags())

	return c
}

func migrateFlagSet(fs *pflag.FlagSet) string {
	dir := fs.StringP("migrations", "m", "./migrations", "migrations directory")

	return *dir
}
