package taskutil

import (
	"golang.org/x/net/proxy"
	"nebula/third_party/http"
	"net/url"
)

func NewCharlesDialer() proxy.ContextDialer {
	u, _ := url.Parse("http://127.0.0.1:8888")
	d, _ := http.NewConnectDialer(u)
	return d
}
