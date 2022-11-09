package http

import (
	"context"
	"net/url"
	"strings"
)

type SimpleClient struct {
	ctx       context.Context
	inner     *Client
	defaultFp *FingerprintConfig
}

func NewSimpleClient(ctx context.Context, clientConfig *ClientConfig, fpConfig *FingerprintConfig) *SimpleClient {
	return &SimpleClient{
		ctx:       ctx,
		inner:     NewClient(clientConfig, fpConfig),
		defaultFp: fpConfig,
	}
}

func (sc *SimpleClient) Builder(url string) *RequestBuilder {
	if sc.defaultFp != nil {
		return sc.inner.Builder(sc.ctx, url).Fingerprint(*sc.defaultFp)
	}
	return sc.inner.Builder(sc.ctx, url)
}

func (sc *SimpleClient) Cookies(u *url.URL) []*Cookie {
	if sc.inner.Jar == nil {
		return nil
	}
	return sc.inner.Jar.Cookies(u)
}

func (sc *SimpleClient) SetCookies(u *url.URL, cookies []*Cookie) {
	if sc.inner.Jar == nil {
		return
	}
	sc.inner.Jar.SetCookies(u, cookies)
}

func (sc *SimpleClient) AddCookie(u *url.URL, cookie *Cookie) {
	if sc.inner.Jar == nil {
		return
	}
	sc.inner.Jar.SetCookies(u, append(sc.inner.Jar.Cookies(u), cookie))
}

func (sc *SimpleClient) DeleteCookie(u *url.URL, name string) {
	if sc.inner.Jar == nil {
		return
	}
	cookies := sc.inner.Jar.Cookies(u)
	for cookieIdx, cookie := range cookies {
		if strings.ToLower(cookie.Name) == strings.ToLower(name) {
			cookies[cookieIdx] = cookies[len(cookies)-1]
			cookies = cookies[:len(cookies)-1]
			break
		}
	}
	sc.inner.Jar.SetCookies(u, cookies)
}

func (sc *SimpleClient) Inner() *Client {
	return sc.inner
}
