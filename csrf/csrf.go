package csrf

import (
	"html/template"
	"net/http"

	gorillacsrf "github.com/gorilla/csrf"
)

type Config struct {
	Key            []byte
	Secure         bool
	TrustedOrigins []string
}

func Protect(config Config, next http.Handler) http.Handler {
	options := []gorillacsrf.Option{
		gorillacsrf.Secure(config.Secure),
		gorillacsrf.SameSite(gorillacsrf.SameSiteLaxMode),
	}

	if len(config.TrustedOrigins) > 0 {
		options = append(
			options,
			gorillacsrf.TrustedOrigins(config.TrustedOrigins),
		)
	}

	middleware := gorillacsrf.Protect(
		config.Key,
		options...,
	)

	return middleware(next)
}

func Token(r *http.Request) string {
	return gorillacsrf.Token(r)
}

func TemplateField(r *http.Request) template.HTML {
	return gorillacsrf.TemplateField(r)
}
