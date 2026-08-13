package csrf

import (
	"html/template"
	"net/http"

	gorillacsrf "github.com/gorilla/csrf"
)

type Config struct {
	Key    []byte
	Secure bool
}

func Protect(config Config, next http.Handler) http.Handler {
	middleware := gorillacsrf.Protect(
		config.Key,
		gorillacsrf.Secure(config.Secure),
		gorillacsrf.SameSite(gorillacsrf.SameSiteLaxMode),
	)

	return middleware(next)
}

func Token(r *http.Request) string {
	return gorillacsrf.Token(r)
}

func TemplateField(r *http.Request) template.HTML {
	return gorillacsrf.TemplateField(r)
}
