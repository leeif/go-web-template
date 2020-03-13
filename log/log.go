package log

import (
	"io"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/leeif/go-web-template/config"
	"go.uber.org/fx"
)

type Log struct {
	logger log.Logger
}

func (pl *Log) Error(message interface{}) {
	level.Error(pl.logger).Log("error", message)
}

func (pl *Log) Info(message interface{}) {
	level.Info(pl.logger).Log("info", message)
}

func (pl *Log) Debug(message interface{}) {
	level.Debug(pl.logger).Log("debug", message)
}

func (pl *Log) Warn(message interface{}) {
	level.Debug(pl.logger).Log("warning", message)
}

func (pl *Log) With(keyvals ...interface{}) *Log {
	l := log.With(pl.logger, keyvals...)
	return &Log{
		logger: l,
	}
}

func NewLogger(lc fx.Lifecycle, config *config.Config) (*Log, error) {
	var l log.Logger
	var writer io.WriteCloser
	var err error
	if *config.Log.File == "" {
		writer = os.Stdout
	} else {
		writer, err = os.OpenFile(*config.Log.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
	}
	if config.Log.Format.String() == "logfmt" {
		l = log.NewLogfmtLogger(log.NewSyncWriter(writer))
	} else {
		l = log.NewJSONLogger(log.NewSyncWriter(writer))
	}
	l = level.NewFilter(l, config.Log.Level.GetLevelOption())
	l = log.With(l, "ts", log.DefaultTimestamp)
	return &Log{
		logger: l,
	}, nil
}
