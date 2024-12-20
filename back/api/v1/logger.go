package v1

import (
	"github.com/savsgio/atreugo/v11"
	"github.com/sirupsen/logrus"
)

func NewLogger(ctx *atreugo.RequestCtx) *logrus.Entry {
	logger := logrus.New()

	return logger.
		WithField("request_method", string(ctx.Method())).
		WithField("request_path", string(ctx.Path())).
		WithField("request_ip", ctx.RemoteIP().String())
}
