package db

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/webitel/wlog"
)

var maxArgLen = 32

type traceQueryCtxKey struct{}

type traceQueryData struct {
	startTime time.Time
	sql       string
	args      []any
}

type tracer struct {
	log *wlog.Logger
}

func newTracer(log *wlog.Logger) *tracer {
	return &tracer{log: log}
}

func (t *tracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	return context.WithValue(ctx, traceQueryCtxKey{}, &traceQueryData{
		startTime: time.Now(),
		sql:       data.SQL,
		args:      data.Args,
	})
}

func (t *tracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	queryData := ctx.Value(traceQueryCtxKey{}).(*traceQueryData)
	endTime := time.Now()
	interval := endTime.Sub(queryData.startTime)
	pgConn := conn.PgConn()

	log := t.log.With(wlog.String("sql", queryData.sql), wlog.Any("args", logQueryArgs(queryData.args)),
		wlog.Any("time", interval))

	if pgConn != nil {
		pid := pgConn.PID()
		if pid != 0 {
			log = log.With(wlog.Int("pid", int(pid)))
		}
	}

	if data.Err != nil {
		log.Error("database query", wlog.Err(data.Err))

		return
	}

	log.Debug("database query", wlog.String("command_tag", data.CommandTag.String()))
}

func logQueryArgs(args []any) []any {
	logArgs := make([]any, 0, len(args))

	for _, a := range args {
		switch v := a.(type) {
		case []byte:
			if len(v) < maxArgLen {
				a = hex.EncodeToString(v)
			} else {
				a = fmt.Sprintf("%x (truncated %d bytes)", v[:maxArgLen], len(v)-maxArgLen)
			}
		case string:
			if len(v) > maxArgLen {
				var l = 0
				for w := 0; l < maxArgLen; l += w {
					_, w = utf8.DecodeRuneInString(v[l:])
				}

				if len(v) > l {
					a = fmt.Sprintf("%s (truncated %d bytes)", v[:l], len(v)-l)
				}
			}
		}
		logArgs = append(logArgs, a)
	}

	return logArgs
}
