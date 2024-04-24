package httpserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/webitel/wlog"
)

func (s *Server) logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func(start time.Time) {
			requestID := w.Header().Get("X-Request-Id")
			if requestID == "" {
				requestID = "unknown"
			}

			s.log.Info("processed request", wlog.String("request_id", requestID), wlog.String("method", req.Method),
				wlog.String("path", req.URL.Path), wlog.String("remote", req.RemoteAddr), wlog.String("ua", req.UserAgent()),
				wlog.Any("duration", time.Since(start)))
		}(time.Now())

		h.ServeHTTP(w, req)
	})
}

func (s *Server) recoverPanic(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				// Acts as a trigger to make HTTP server automatically close the current connection after a response has been sent.
				w.Header().Set("Connection", "close")
				JSON(w, Error{Msg: fmt.Sprintf("%s", err)}, http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	})
}

func (s *Server) authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("access_token")
		if len(token) == 0 {
			JSON(w, Error{Msg: "access token is missing"}, http.StatusUnauthorized)

			return
		}

		if token != s.cfg.Server.Token {
			JSON(w, Error{Msg: "unauthorized"}, http.StatusUnauthorized)

			return
		}

		h.ServeHTTP(w, r)
	})
}
