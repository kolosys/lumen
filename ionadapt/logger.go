package ionadapt

import (
	"fmt"

	"github.com/kolosys/lumen/logs"
)

type logger struct{ l *logs.Logger }

func (a logger) Debug(msg string, kv ...any) { a.l.Debug(msg, logFields(kv)...) }
func (a logger) Info(msg string, kv ...any)  { a.l.Info(msg, logFields(kv)...) }
func (a logger) Warn(msg string, kv ...any)  { a.l.Warn(msg, logFields(kv)...) }

func (a logger) Error(msg string, err error, kv ...any) {
	fs := logFields(kv)
	a.l.Error(msg, append(fs, logs.Err(err))...)
}

func logFields(kv []any) []logs.Field {
	n := len(kv) / 2
	if n == 0 {
		return nil
	}
	fs := make([]logs.Field, 0, n)
	eachKV(kv, func(key string, val any) {
		fs = append(fs, logs.String(key, fmt.Sprint(val)))
	})
	return fs
}
