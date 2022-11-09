package proxies

import "net/url"

type ProxyGroup struct {
	Name    string    `json:"name"`
	Proxies *proxySet `json:"proxies"`
}

type ProxyData struct {
	Id  string   `json:"id"`
	Url *url.URL `json:"url"`
}
