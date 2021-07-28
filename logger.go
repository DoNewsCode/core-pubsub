package pubsub

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Logger implements watermill.Logger
type Logger struct {
	next log.Logger
}

// NewLogger constructs a Logger
func NewLogger(logger log.Logger) Logger {
	return Logger{next: logger}
}

// Error implements watermill.Logger
func (l Logger) Error(msg string, err error, fields watermill.LogFields) {
	args := []interface{}{"msg", msg, "err", err}
	for k, v := range fields {
		args = append(args, k, v)
	}
	level.Error(l.next).Log(args...)
}

// Info implements watermill.Logger
func (l Logger) Info(msg string, fields watermill.LogFields) {
	args := []interface{}{"msg", msg}
	for k, v := range fields {
		args = append(args, k, v)
	}
	level.Info(l.next).Log(args...)
}

// Debug implements watermill.Logger
func (l Logger) Debug(msg string, fields watermill.LogFields) {
	args := []interface{}{"msg", msg}
	for k, v := range fields {
		args = append(args, k, v)
	}
	level.Debug(l.next).Log(args...)
}

// Trace implements watermill.Logger
func (l Logger) Trace(msg string, fields watermill.LogFields) {
	args := []interface{}{"msg", msg}
	for k, v := range fields {
		args = append(args, k, v)
	}
	level.Debug(l.next).Log(args...)
}

// With implements watermill.Logger
func (l Logger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	var args []interface{}
	for k, v := range fields {
		args = append(args, k, v)
	}
	return Logger{next: log.With(l.next, args...)}
}

