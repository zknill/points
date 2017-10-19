package log

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/appengine/log"
)

// Logger methods for logging in the points package.
// It allows to flip between logging to file / stdout
// and logging to a different location like app engine.
type Logger interface {
	Infof(ctx context.Context, format string, args ...interface{})
	Warningf(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
}

var _ Logger = (*AppEngine)(nil)
var _ Logger = (*Logrus)(nil)

// DefaultLogger creates a logrus logger that will write
// log messages to stdout, with the default logrus formatter.
func DefaultLogger() Logger {
	return &Logrus{log: logrus.New()}
}

// AppEngine is a logger that sends log messages to google app engine log
type AppEngine struct{}

// Infof log to app engine log
func (*AppEngine) Infof(ctx context.Context, format string, args ...interface{}) {
	log.Infof(ctx, format, args...)
}

// Warningf log to app engine log
func (*AppEngine) Warningf(ctx context.Context, format string, args ...interface{}) {
	log.Warningf(ctx, format, args...)
}

// Errorf log to app engine log
func (*AppEngine) Errorf(ctx context.Context, format string, args ...interface{}) {
	log.Errorf(ctx, format, args...)
}

// Logrus wraps a logrus logger and implements the Logger interface
type Logrus struct {
	log *logrus.Logger
}

// Infof log to logrus logger
func (ll *Logrus) Infof(_ context.Context, format string, args ...interface{}) {
	ll.log.Infof(format, args...)
}

// Warningf log to logrus logger
func (ll *Logrus) Warningf(_ context.Context, format string, args ...interface{}) {
	ll.log.Warningf(format, args...)
}

// Errorf log to logrus logger
func (ll *Logrus) Errorf(_ context.Context, format string, args ...interface{}) {
	ll.log.Errorf(format, args...)
}
