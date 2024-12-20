package v1

import "github.com/savsgio/atreugo/v11"

func ping(ctx *atreugo.RequestCtx) error {
	return ctx.JSONResponse(map[string]string{"message": "Pong!"})
}
