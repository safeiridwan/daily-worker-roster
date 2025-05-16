package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var ErrRecordNotFound = errors.New("record not found")

type LogLevel = gormlogger.LogLevel

const (
	Silent = gormlogger.Silent
	Error  = gormlogger.Error
	Warn   = gormlogger.Warn
	Info   = gormlogger.Info
)

type Config struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	LogLevel                  LogLevel
}

type Interface interface {
	LogMode(LogLevel) gormlogger.Interface
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}

func New(writer zerolog.Logger, config Config) Interface {
	return &logger{
		writer: writer,
		Config: config,
	}
}

type logger struct {
	writer zerolog.Logger
	Config
}

func (l *logger) LogMode(level LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *logger) log(ctx context.Context, level LogLevel, msg string, data ...interface{}) {
	if l.LogLevel >= level {
		event := l.writer.WithLevel(zerolog.InfoLevel)
		if id, ok := hlog.IDFromCtx(ctx); ok {
			event = event.Str("request_id", id.String())
		}
		event.Msgf("%s %s", utils.FileWithLineNum(), fmt.Sprintf(msg, data...))
	}
}

func (l *logger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.log(ctx, Info, msg, data...)
}

func (l *logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.log(ctx, Warn, msg, data...)
}

func (l *logger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.log(ctx, Error, msg, data...)
}

func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.LogLevel >= Error && (!errors.Is(err, ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		l.log(ctx, Error, "[%.3fms] [rows:%v] %s %v", float64(elapsed.Nanoseconds())/1e6, rows, sql, err)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= Warn:
		l.log(ctx, Warn, "[%.3fms] [rows:%v] %s %s", float64(elapsed.Nanoseconds())/1e6, rows, sql, fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold))
	case l.LogLevel == Info:
		l.log(ctx, Info, "[%.3fms] [rows:%v] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
}

func (l *logger) ParamsFilter(_ context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.Config.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}

type traceRecorder struct {
	Interface
	BeginAt      time.Time
	SQL          string
	RowsAffected int64
	Err          error
}

func (l *traceRecorder) New() *traceRecorder {
	return &traceRecorder{Interface: l.Interface, BeginAt: time.Now()}
}

func (l *traceRecorder) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	l.BeginAt = begin
	l.SQL, l.RowsAffected = fc()
	l.Err = err
}
