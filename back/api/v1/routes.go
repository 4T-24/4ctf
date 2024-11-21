package v1

import (
	"github.com/savsgio/atreugo/v11"
)

func SetupRoutes(router *atreugo.Atreugo) {
	group := router.NewGroupPath("/api/v1")

	group.GET("/ping", ping)
}
