package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/webitel/wlog"

	"github.com/kirychukyurii/fd-import/config"
)

type Server struct {
	cfg *config.Config
	log *wlog.Logger

	srv    *http.Server
	router *http.ServeMux
}

func New(cfg *config.Config, log *wlog.Logger) *Server {
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:     cfg.Server.Address,
		ErrorLog: log.StdLog(),
	}

	s := &Server{
		cfg:    cfg,
		log:    log,
		srv:    srv,
		router: mux,
	}

	s.srv.Handler = s.recoverPanic(s.logging(s.authenticate(s.router)))

	return s
}

func (s *Server) Serve() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
