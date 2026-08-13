package csrf

import "net/http"

type Config struct {
	TrustedOrigins []string
}

func Protect(
	config Config,
	next http.Handler,
) (http.Handler, error) {
	protection := http.NewCrossOriginProtection()

	for _, origin := range config.TrustedOrigins {
		if err := protection.AddTrustedOrigin(origin); err != nil {
			return nil, err
		}
	}

	return protection.Handler(next), nil
}
