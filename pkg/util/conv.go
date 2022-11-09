package util

import (
	"nebula/third_party/http"

	"github.com/Nebulabots/go-astilectron"
)

func ConvertAstilectronCookies(cookies []astilectron.SessionCookie) (httpCookies []*http.Cookie) {
	httpCookies = make([]*http.Cookie, 0)
	for _, c := range cookies {
		expires := ""
		if c.ExpirationDate != nil {
			expires = string(int(*c.ExpirationDate))
		}

		n := http.Cookie{
			Name:       c.Name,
			Value:      c.Value,
			Path:       c.Path,
			Domain:     c.Domain,
			RawExpires: expires,
			Secure:     *c.Secure,
			HttpOnly:   *c.HttpOnly,
		}

		httpCookies = append(httpCookies, &n)
	}

	return
}
