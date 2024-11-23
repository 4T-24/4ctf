package v1

import (
	"encoding/json"

	"github.com/savsgio/atreugo/v11"
	"github.com/sirupsen/logrus"
	"gopkg.in/validator.v2"
)

func WithDefaults(api *Api, fn func(api *Api) func(ctx *atreugo.RequestCtx) error) func(ctx *atreugo.RequestCtx) error {
	return func(ctx *atreugo.RequestCtx) error {
		session, _ := api.session.GetSession(ctx.RequestCtx)
		sessionBefore := *session

		ctx.SetUserValue("session", session)

		api = &Api{
			config:  api.config,
			session: api.session,
			Entry:   NewLogger(ctx),
		}

		err := fn(api)(ctx)
		if err != nil {
			return err
		}

		// Check if the session has been updated
		if sessionBefore != *session {
			return api.session.SetSession(ctx.RequestCtx, session)
		}
		return nil
	}
}

func WithBody[K any](fn func(api *Api) func(ctx *atreugo.RequestCtx, body K) error) func(api *Api) func(ctx *atreugo.RequestCtx) error {
	return func(api *Api) func(ctx *atreugo.RequestCtx) error {
		return func(ctx *atreugo.RequestCtx) error {
			var body K
			if err := json.Unmarshal(ctx.PostBody(), &body); err != nil {
				logrus.
					WithField("request_ip", ctx.RemoteIP().String()).
					WithError(err).
					Warn("bad request")
				return ctx.JSONResponse(NewErrorResponse(400, []ResponseError{{Message: "Bad request"}}))
			}

			if err := validator.Validate(body); err != nil {
				logrus.
					WithField("request_ip", ctx.RemoteIP().String()).
					WithError(err).
					Warn("invalid request")
				return ctx.JSONResponse(NewErrorResponse(400, validatorErrorToResponseError(err)))
			}

			return fn(api)(ctx, body)
		}
	}
}

func WithValidSession(fn func(api *Api) func(ctx *atreugo.RequestCtx) error) func(api *Api) func(ctx *atreugo.RequestCtx) error {
	return func(api *Api) func(ctx *atreugo.RequestCtx) error {
		return func(ctx *atreugo.RequestCtx) error {
			session := ctx.UserValue("session").(*Session)
			if !session.Valid {
				return ctx.JSONResponse(NewErrorResponse(401, []ResponseError{{Message: "Unauthorized"}}))
			}

			return fn(api)(ctx)
		}
	}
}
