package v1

import (
	"encoding/json"
	"net/http"

	"github.com/savsgio/atreugo/v11"
	"github.com/sirupsen/logrus"
	"gopkg.in/validator.v2"
)

func Default(api *Api, fn func(api *Api) func(ctx *atreugo.RequestCtx) *Response[any]) func(ctx *atreugo.RequestCtx) error {
	return func(ctx *atreugo.RequestCtx) error {
		session, _ := api.session.GetSession(ctx.RequestCtx)
		sessionBefore := *session

		ctx.SetUserValue("session", session)

		api = &Api{
			config:  api.config,
			session: api.session,
			Entry:   NewLogger(ctx),
		}

		// Log the request
		entry := api.WithField("request_size", ctx.Request.Header.ContentLength())
		if len(string(ctx.URI().QueryString())) > 0 {
			entry = entry.WithField("request_query", string(ctx.URI().QueryString()))
		}
		entry.Info("request received")

		response := fn(api)(ctx)

		// Check if the session has been updated
		if sessionBefore != *session {
			err := api.session.SetSession(ctx.RequestCtx, session)
			if err != nil {
				api.
					WithError(err).
					Error("cannot update session")
			}
		}

		return response.Send(ctx)
	}
}

func Body[K any](fn func(api *Api) func(ctx *atreugo.RequestCtx, body K) *Response[any]) func(api *Api) func(ctx *atreugo.RequestCtx) *Response[any] {
	return func(api *Api) func(ctx *atreugo.RequestCtx) *Response[any] {
		return func(ctx *atreugo.RequestCtx) *Response[any] {
			var body K
			if err := json.Unmarshal(ctx.PostBody(), &body); err != nil {
				logrus.
					WithField("request_ip", ctx.RemoteIP().String()).
					WithError(err).
					Warn("bad request")
				return Error(http.StatusBadRequest, "Bad request")
			}

			if err := validator.Validate(body); err != nil {
				logrus.
					WithField("request_ip", ctx.RemoteIP().String()).
					WithError(err).
					Warn("invalid request")
				return Errors(http.StatusBadRequest, validatorErrorToResponseError(err))
			}

			return fn(api)(ctx, body)
		}
	}
}

func RequireValidSession(fn func(api *Api) func(ctx *atreugo.RequestCtx) *Response[any]) func(api *Api) func(ctx *atreugo.RequestCtx) *Response[any] {
	return func(api *Api) func(ctx *atreugo.RequestCtx) *Response[any] {
		return func(ctx *atreugo.RequestCtx) *Response[any] {
			session := ctx.UserValue("session").(*Session)
			if !session.Valid {
				return Error(http.StatusUnauthorized, "Unauthorized")
			}

			return fn(api)(ctx)
		}
	}
}
