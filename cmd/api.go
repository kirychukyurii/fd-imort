package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/webitel/wlog"

	"github.com/kirychukyurii/fd-import/config"
	"github.com/kirychukyurii/fd-import/pkg/db"
	"github.com/kirychukyurii/fd-import/pkg/httpserver"
)

func apiCommand(cfg *config.Config, log *wlog.Logger) *cobra.Command {
	c := &cobra.Command{
		Use:          "api",
		Short:        "API for download ticket attachments",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			dbpool, err := db.New(cmd.Context(), log, cfg.DSN)
			if err != nil {
				return err
			}

			srv := httpserver.New(cfg, log)
			srv.RegisterHandlers(dbpool)
			a := api{
				cfg:    cfg,
				log:    log,
				srv:    srv,
				dbpool: dbpool,
			}

			if err := a.run(cmd.Context()); err != nil {
				a.log.Error("run api", wlog.Err(err))
			}

			return err
		},
	}

	apiFlags(c.PersistentFlags(), cfg)

	return c
}

func apiFlags(fs *pflag.FlagSet, cfg *config.Config) {
	fs.StringVarP(&cfg.Server.Address, "bind", "b", "0.0.0.0:10111", "bind address")
	fs.StringVarP(&cfg.Server.Token, "access-token", "t", "", "access token")
}

type api struct {
	cfg *config.Config
	log *wlog.Logger

	srv    *httpserver.Server
	dbpool *db.Connection

	errorCh chan error
}

func (a *api) run(ctx context.Context) error {
	a.errorCh = make(chan error)
	defer close(a.errorCh)
	go func() {
		a.log.Info("listening http", wlog.String("server", a.cfg.Server.Address))
		if err := a.srv.Serve(); err != nil {
			a.errorCh <- err
		}
	}()

	// api blocks until it receives a signal to exit
	// this signal may come from the node or from sig-abort (ctrl-c)
	select {
	case <-ctx.Done():
		return nil
	case err := <-a.errorCh:
		return err
	}
}
