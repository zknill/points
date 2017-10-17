package points

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/appengine/log"
)

type logger interface {
	Infof(ctx context.Context, format string, args ...interface{})
	Warningf(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
}

func defaultLogger() logger {
	return &logrusLog{log: logrus.New()}
}

type appEngineLog struct{}

func (*appEngineLog) Infof(ctx context.Context, format string, args ...interface{}) {
	log.Infof(ctx, format, args...)
}

func (*appEngineLog) Warningf(ctx context.Context, format string, args ...interface{}) {
	log.Warningf(ctx, format, args...)
}

func (*appEngineLog) Errorf(ctx context.Context, format string, args ...interface{}) {
	log.Errorf(ctx, format, args...)
}

type logrusLog struct {
	log *logrus.Logger
}

func (ll *logrusLog) Infof(_ context.Context, format string, args ...interface{}) {
	ll.log.Infof(format, args...)
}

func (ll *logrusLog) Warningf(_ context.Context, format string, args ...interface{}) {
	ll.log.Warningf(format, args...)
}

func (ll *logrusLog) Errorf(_ context.Context, format string, args ...interface{}) {
	ll.log.Errorf(format, args...)
}
