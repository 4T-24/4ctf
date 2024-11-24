package v1

import (
	"4ctf/config"
	"4ctf/translations"
	"4ctf/utils"

	"github.com/nicksnyder/go-i18n/v2/i18n"
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

	localizer *i18n.Localizer

	*logrus.Entry
}

func (api *Api) Translate(messageID string) string {
	if api.localizer == nil {
		return messageID
	}
	return api.localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: messageID + ".message"})
}

func SetupRoutes(router *atreugo.Atreugo, config *config.Config) {
	var api = &Api{
		config:    config,
		session:   utils.NewSessionManager([]byte(config.Server.Key), nil, "session", Session{}),
		localizer: i18n.NewLocalizer(translations.Bundle, "en"),
	}
	group := router.NewGroupPath("/api/v1")

	// Validators
	setupValidators()

	// Routes
	setupAuthRoutes(api, group)

	// Misc
	group.GET("/ping", ping)
}
