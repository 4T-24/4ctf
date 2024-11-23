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

func SetupRoutes(router *atreugo.Atreugo, config *config.Config) {
	var api = &Api{
		config:  config,
		session: utils.NewSessionManager([]byte(config.Server.Key), nil, "session", Session{}),
	}
	group := router.NewGroupPath("/api/v1")

	// Validators
	setupValidators()

	// Routes
	setupAuthRoutes(api, group)

	// Misc
	group.GET("/ping", ping)
}
