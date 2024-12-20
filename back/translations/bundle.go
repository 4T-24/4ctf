package translations

import (
	"embed"
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed active.*.json
var LocaleFS embed.FS

var Bundle *i18n.Bundle

func init() {
	bundle := i18n.NewBundle(language.English)

	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.LoadMessageFileFS(LocaleFS, "active.en.json")
	bundle.LoadMessageFileFS(LocaleFS, "active.fr.json")

	Bundle = bundle
}
