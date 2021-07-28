package pubsub

import (
	"bytes"
	"testing"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	cases := []struct {
		name    string
		f       func(logger Logger)
		asserts func(buf string, t *testing.T)
	}{
		{
			"debug",
			func(logger Logger) {
				logger.Debug("foo", nil)
			},
			func(buf string, t *testing.T) {
				assert.Contains(t, buf, "level=debug")
			},
		},
		{
			"trace",
			func(logger Logger) {
				logger.Debug("foo", nil)
			},
			func(buf string, t *testing.T) {
				assert.Contains(t, buf, "level=debug")
			},
		},
		{
			"info",
			func(logger Logger) {
				logger.Info("foo", nil)
			},
			func(buf string, t *testing.T) {
				assert.Contains(t, buf, "level=info")
			},
		},
		{
			"error",
			func(logger Logger) {
				logger.Error("foo", nil, nil)
			},
			func(buf string, t *testing.T) {
				assert.Contains(t, buf, "level=error")
			},
		},
		{
			"with",
			func(logger Logger) {
				logger.With(watermill.LogFields{"foo": "bar"}).Error("baz", nil, nil)
			},
			func(buf string, t *testing.T) {
				assert.Contains(t, buf, "foo=bar")
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger(log.NewSyncLogger(log.NewLogfmtLogger(&buf)))
			c.f(logger)
			c.asserts(buf.String(), t)
		})
	}

}
