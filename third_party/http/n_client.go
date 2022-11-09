package http

import (
	"bytes"
	"context"
	"encoding/json"
	"golang.org/x/net/proxy"
	"io/ioutil"
	tls "nebula/third_party/utls"
	"net/url"
	"strings"
	"time"
)

var (
	DefaultClientConfig = &ClientConfig{
		ClientHello: tls.HelloChrome_Auto,
		TlsConfig:   &tls.Config{NextProtos: []string{Http2NextProtoTLS, "http/1.1"}},
		ProxyDialer: proxy.Direct,
	}

	InsecureClientConfig = &ClientConfig{
		ClientHello: tls.HelloChrome_Auto,
		TlsConfig:   &tls.Config{InsecureSkipVerify: true, NextProtos: []string{Http2NextProtoTLS, "http/1.1"}},
		ProxyDialer: proxy.Direct,
	}
)

type ClientConfig struct {
	ClientHello tls.ClientHelloID
	TlsConfig   *tls.Config

	ProxyDialer proxy.ContextDialer
}

func (c *ClientConfig) Clone() *ClientConfig {
	return &ClientConfig{
		ClientHello: c.ClientHello,
		TlsConfig:   c.TlsConfig.Clone(),
		ProxyDialer: c.ProxyDialer,
	}
}

func NewClient(clientConfig *ClientConfig, fpConfig *FingerprintConfig) *Client {
	return &Client{
		Timeout:   time.Second * 30,
		Transport: NewTlsTransport(clientConfig, fpConfig),
	}
}

type RequestBuilder struct {
	sender *Client
	inner  *Request
	err    error
}

func newRequestBuilder(ctx context.Context, sender *Client, url string) *RequestBuilder {
	inner, err := NewRequestWithContext(ctx, MethodGet, url, nil)
	return &RequestBuilder{inner: inner, sender: sender, err: err}
}

func (c *Client) Builder(ctx context.Context, url string) *RequestBuilder {
	return newRequestBuilder(ctx, c, url)
}

func (r *RequestBuilder) Method(method string) *RequestBuilder {
	if r.err != nil {
		return r
	}

	r.inner.Method = method
	return r
}

func (r *RequestBuilder) Fingerprint(config FingerprintConfig) *RequestBuilder {
	r.inner.FingerprintConfig = config
	return r
}

func (r *RequestBuilder) addOrderedHeader(key string) {
	if headerOrder, ok := r.inner.Header[HeaderOrderKey]; ok {
		r.inner.Header[HeaderOrderKey] = append(headerOrder, key)
	} else {
		r.inner.Header[HeaderOrderKey] = []string{key}
	}
}

func (r *RequestBuilder) addHeader(key, value string, ordered bool) *RequestBuilder {
	if r.err != nil {
		return r
	}

	if header, ok := r.inner.Header[key]; ok {
		header = append(header, value)
		r.inner.Header[key] = header
	} else {
		if ordered {
			r.addOrderedHeader(key)
		}
		r.inner.Header[key] = []string{value}
	}
	return r
}

func (r *RequestBuilder) Header(key, value string) *RequestBuilder {
	return r.addHeader(key, value, true)
}

func (r *RequestBuilder) HeaderUnordered(key, value string) *RequestBuilder {
	return r.addHeader(key, value, false)
}

func (r *RequestBuilder) OverwriteHeader(key, value string) *RequestBuilder {
	if r.err != nil {
		return r
	}

	r.inner.Header[key] = []string{value}
	r.addOrderedHeader(key)
	return r
}

func (r *RequestBuilder) HeaderOrder(keys ...string) *RequestBuilder {
	r.inner.Header[HeaderOrderKey] = keys
	return r
}

func (r *RequestBuilder) Body(body []byte) *RequestBuilder {
	r.inner.Body = ioutil.NopCloser(bytes.NewReader(body))
	return r
}

func (r *RequestBuilder) FormBody(formBody url.Values) *RequestBuilder {
	return r.StringBody(formBody.Encode())
}

func (r *RequestBuilder) StringBody(body string) *RequestBuilder {
	r.inner.Body = ioutil.NopCloser(strings.NewReader(body))
	return r
}

func (r *RequestBuilder) JsonBody(obj interface{}) *RequestBuilder {
	if r.err != nil {
		return r
	}

	jsonBody, err := json.Marshal(obj)
	if err != nil {
		r.err = err
		return r
	}

	r.Body(jsonBody)
	return r
}

func (r *RequestBuilder) Send() (*Response, error) {
	if r.err != nil {
		return nil, r.err
	}

	return r.sender.Do(r.inner)
}

func (r *RequestBuilder) SendAndClose() (*Response, error) {
	resp, err := r.Send()
	if err != nil {
		return nil, err
	}

	respBody := resp.readAndClose()
	resp.cachedBody = respBody
	return resp, nil
}
