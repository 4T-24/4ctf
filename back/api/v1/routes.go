package v1

import (
	"4ctf/config"
	"4ctf/utils"

	"github.com/savsgio/atreugo/v11"
	"github.com/sirupsen/logrus"
)

type Session struct {
	Valid         bool
	UserSessionID uint64
}

type Api struct {
	config  *config.Config
	session *utils.SessionManager[Session]

	*logrus.Entry
}

func SetupRoutes(router *atreugo.Atreugo, cfg *config.Config) {
	var api = &Api{
		config:  cfg,
		session: utils.NewSessionManager([]byte(cfg.Server.Key), nil, "session", Session{}),
	}

	// In dev, we want CORS to allow everyone
	if cfg.IsDevelopment() {
		router.UseFinal(func(ctx *atreugo.RequestCtx) {
			ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
			ctx.Response.Header.SetBytesV("Access-Control-Allow-Origin", ctx.Request.Header.Peek("Origin"))
		})
	}

	group := router.NewGroupPath("/api/v1")

	// Validators
	setupValidators()

	// Routes
	setupAuthRoutes(api, group)

	// Misc
	group.GET("/ping", ping)
}
