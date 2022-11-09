package ys

import (
	"nebula/third_party/http"
	"net/url"
	"sync"
)

type SingleJar struct {
	cookies  map[string]*http.Cookie
	cookieMu sync.Mutex
}

func newSingleJar() *SingleJar {
	return &SingleJar{
		cookies: make(map[string]*http.Cookie),
	}
}

func (j *SingleJar) SetCookies(_ *url.URL, cookies []*http.Cookie) {
	j.cookieMu.Lock()
	for _, cookie := range cookies {
		j.cookies[cookie.Name] = cookie
	}
	j.cookieMu.Unlock()
}

func (j *SingleJar) Cookies(_ *url.URL) []*http.Cookie {
	j.cookieMu.Lock()
	var cookies []*http.Cookie
	for _, cookie := range j.cookies {
		cookies = append(cookies, cookie)
	}
	j.cookieMu.Unlock()
	return cookies
}

func (j *SingleJar) DeleteCookie(name string) {
	j.cookieMu.Lock()
	delete(j.cookies, name)
	j.cookieMu.Unlock()
}
