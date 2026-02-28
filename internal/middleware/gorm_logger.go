package middleware

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

type Logger struct {
	zap            *zap.Logger
	level          glogger.LogLevel
	slow           time.Duration
	ignoreNotFound bool
}

func NewGormLogger(z *zap.Logger, level glogger.LogLevel, slow time.Duration, ignoreNotFound bool) *Logger {
	return &Logger{
		zap:            z,
		level:          level,
		slow:           slow,
		ignoreNotFound: ignoreNotFound,
	}
}

func (l *Logger) LogMode(level glogger.LogLevel) glogger.Interface {
	n := *l
	n.level = level
	return &n
}

func (l *Logger) Info(ctx context.Context, msg string, data ...any) {
	if l.level < glogger.Info {
		return
	}
	l.zap.Sugar().Infow(msg, "data", data)
}

func (l *Logger) Warn(ctx context.Context, msg string, data ...any) {
	if l.level < glogger.Warn {
		return
	}
	l.zap.Sugar().Warnw(msg, "data", data)
}

func (l *Logger) Error(ctx context.Context, msg string, data ...any) {
	if l.level < glogger.Error {
		return
	}
	l.zap.Sugar().Errorw(msg, "data", data)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.level == glogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 忽略 record not found（可选）
	if l.ignoreNotFound && errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}

	fields := []zap.Field{
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
		zap.String("sql", sql),
	}

	switch {
	case err != nil && l.level >= glogger.Error:
		l.zap.Error("gorm query", append(fields, zap.Error(err))...)
	case l.slow > 0 && elapsed > l.slow && l.level >= glogger.Warn:
		l.zap.Warn("gorm slow query", fields...)
	case l.level >= glogger.Info:
		l.zap.Info("gorm query", fields...)
	}
}
