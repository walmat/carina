package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	stdhttp "net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"

	"golang.org/x/net/proxy"

	utls "nebula/third_party/utls"
)

var errProtocolNegotiated = errors.New("protocol negotiated")

type tlsRoundTripper struct {
	sync.Mutex

	clientTlsConfig *utls.Config
	clientHelloId   utls.ClientHelloID

	cachedConnections map[string]net.Conn
	cachedTransports  map[string]RoundTripper

	h3cache map[string]map[string]altSvcData

	dialer proxy.ContextDialer

	fpConfig *FingerprintConfig
}

// splitAtCommas split s at commas, ignoring commas in strings.
func splitAtCommas(s string) []string {
	var res []string
	var beg int
	var inString bool

	for i := 0; i < len(s); i++ {
		if s[i] == ',' && !inString {
			res = append(res, s[beg:i])
			beg = i + 1
		} else if s[i] == '"' {
			if !inString {
				inString = true
			} else if i > 0 && s[i-1] != '\\' {
				inString = false
			}
		}
	}
	return append(res, s[beg:])
}

type altSvcData struct {
	Authority string
	MaxAge    *int
}

func (rt *tlsRoundTripper) RoundTrip(req *Request) (*Response, error) {
	addr := rt.getDialTLSAddr(req)
	if _, ok := rt.cachedTransports[addr]; !ok {
		if err := rt.getTransport(req, addr); err != nil {
			return nil, err
		}
	}

	resp, err := rt.cachedTransports[addr].RoundTrip(req)
	if err != nil {
		return resp, err
	}

	if rt.h3cache != nil {
		// TODO: write proper parser that does not rely on splitting and reusing data
		altSvc := resp.Header.Get("Alt-Svc")
		if altSvc != "" {
			if host, _, err := net.SplitHostPort(addr); err == nil {
				altSvcs := splitAtCommas(strings.ReplaceAll(altSvc, " ", ""))
				for _, svc := range altSvcs {
					proto := ""
					parsedSvcData := altSvcData{}
					rawSvcData := strings.Split(svc, ";")
					if len(rawSvcData) >= 1 {
						typeData := strings.Split(rawSvcData[0], "=")
						if len(typeData) == 2 {
							proto = typeData[0]
							parsedSvcData.Authority = strings.ReplaceAll(typeData[1], "\"", "")
						}
					}
					if len(rawSvcData) >= 2 {
						ageData := strings.Split(rawSvcData[1], "=")
						if len(ageData) == 2 {
							if maxAge, err := strconv.Atoi(ageData[1]); err == nil {
								parsedSvcData.MaxAge = &maxAge
							}
						}
					}
					if proto != "" {
						if c, ok := rt.h3cache[host]; ok {
							c[proto] = parsedSvcData
						} else {
							rt.h3cache[host] = make(map[string]altSvcData)
							rt.h3cache[host][proto] = parsedSvcData
						}
					}
				}

				if rt.h3cache[host] != nil {
					quicConfig := &quic.Config{
						MaxIncomingStreams: -1, // don't allow the server to create bidirectional streams
						KeepAlive:          true,
					}
					for proto := range rt.h3cache[host] {
						switch proto {
						case "h3":
							quicConfig.Versions = []quic.VersionNumber{quic.Version1}
							break
						case "h3-34":
							quicConfig.Versions = []quic.VersionNumber{quic.VersionDraft34}
							break
						case "h3-32":
							quicConfig.Versions = []quic.VersionNumber{quic.VersionDraft32}
							break
						case "h3-29":
							quicConfig.Versions = []quic.VersionNumber{quic.VersionDraft29}
							break
						}
					}
					if len(quicConfig.Versions) == 1 {
						rt.cachedTransports[addr] = NewDecompressTransport(&http3Transport{inner: &http3.RoundTripper{QuicConfig: quicConfig, DisableCompression: true}})
					}
				}
			}
		}
	}

	return resp, err
}

func (rt *tlsRoundTripper) getTransport(req *Request, addr string) error {
	switch strings.ToLower(req.URL.Scheme) {
	case "http":
		rt.cachedTransports[addr] = NewDecompressTransport(&Transport{DialTLSContext: rt.dialTLS, DisableCompression: true})
		return nil
	case "https":
	default:
		return fmt.Errorf("invalid URL scheme: [%v]", req.URL.Scheme)
	}

	_, err := rt.dialTLS(context.Background(), "tcp", addr)
	switch err {
	case errProtocolNegotiated:
	case nil:
		rt.Lock()
		defer rt.Unlock()
		if _, ok := rt.cachedTransports[addr]; !ok {
			// Should never happen.
			panic("dialTLS returned no error when determining cachedTransports")
		}
	default:
		return err
	}

	return nil
}

type ignorePushHandler struct{}

func (*ignorePushHandler) HandlePush(r *Http2PushedRequest) {
	pushedResp, err := r.ReadResponse(r.Promise.Context())
	if err == nil {
		_ = pushedResp.Body.Close()
	}
}

func (rt *tlsRoundTripper) dialTLS(ctx context.Context, network, addr string) (net.Conn, error) {
	rt.Lock()
	defer rt.Unlock()

	// If we have the connection from when we determined the HTTPS
	// cachedTransports to use, return that.
	if conn := rt.cachedConnections[addr]; conn != nil {
		delete(rt.cachedConnections, addr)
		return conn, nil
	}

	rawConn, err := rt.dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	var host string
	if host, _, err = net.SplitHostPort(addr); err != nil {
		host = addr
	}

	tlsConfig := rt.clientTlsConfig

	// Set the server name if necessary
	if !tlsConfig.InsecureSkipVerify && tlsConfig.ServerName == "" {
		newConfig := rt.clientTlsConfig.Clone()
		newConfig.ServerName = host
		tlsConfig = newConfig
	}

	conn := utls.UClient(rawConn, tlsConfig, rt.clientHelloId)
	if err = conn.Handshake(); err != nil {
		_ = conn.Close()
		return nil, err
	}

	if rt.cachedTransports[addr] != nil {
		return conn, nil
	}

	// No http.Transport constructed yet, create one based on the results
	// of ALPN.
	switch conn.ConnectionState().NegotiatedProtocol {
	case Http2NextProtoTLS:
		// The remote peer is speaking HTTP 2 + TLS.
		rt.cachedTransports[addr] = NewDecompressTransport(&http2Transport{DialTLS: rt.dialTLSHTTP2, PushHandler: &ignorePushHandler{}, DisableCompression: true, fpConfig: rt.fpConfig})
	default:
		// Assume the remote peer is speaking HTTP 1.x + TLS.
		rt.cachedTransports[addr] = NewDecompressTransport(&Transport{DialTLSContext: rt.dialTLS, DisableCompression: true})
	}

	// Stash the connection just established for use servicing the
	// actual request (should be near-immediate).
	rt.cachedConnections[addr] = conn

	return nil, errProtocolNegotiated
}

func (rt *tlsRoundTripper) dialTLSHTTP2(network, addr string, _ *tls.Config) (net.Conn, error) {
	return rt.dialTLS(context.Background(), network, addr)
}

func (rt *tlsRoundTripper) getDialTLSAddr(req *Request) string {
	host, port, err := net.SplitHostPort(req.URL.Host)
	if err == nil {
		return net.JoinHostPort(host, port)
	}
	return net.JoinHostPort(req.URL.Host, "443") // we can assume port is 443 at this point
}

func NewTlsTransport(clientConfig *ClientConfig, fpConfig *FingerprintConfig) RoundTripper {
	if clientConfig == nil {
		clientConfig = DefaultClientConfig
	}
	if fpConfig == nil {
		fpConfig = DefaultFingerprint
	}
	return &tlsRoundTripper{
		dialer: clientConfig.ProxyDialer,

		clientTlsConfig: clientConfig.TlsConfig,
		clientHelloId:   clientConfig.ClientHello,

		cachedTransports:  make(map[string]RoundTripper),
		cachedConnections: make(map[string]net.Conn),

		h3cache: nil,

		fpConfig: fpConfig,
	}
}

type http3Transport struct {
	inner *http3.RoundTripper
}

func (t *http3Transport) RoundTrip(req *Request) (*Response, error) {
	stdReq := &stdhttp.Request{
		Method:           req.Method,
		URL:              req.URL,
		Proto:            req.Proto,
		ProtoMajor:       req.ProtoMajor,
		ProtoMinor:       req.ProtoMinor,
		Header:           stdhttp.Header(req.Header),
		Body:             req.Body,
		GetBody:          req.GetBody,
		ContentLength:    req.ContentLength,
		TransferEncoding: req.TransferEncoding,
		Close:            req.Close,
		Host:             req.Host,
		Form:             req.Form,
		PostForm:         req.PostForm,
		MultipartForm:    req.MultipartForm,
		Trailer:          stdhttp.Header(req.Trailer),
		RemoteAddr:       req.RemoteAddr,
		RequestURI:       req.RequestURI,
		TLS:              req.TLS,
		Cancel:           req.Cancel,
	}
	stdResp, err := t.inner.RoundTrip(stdReq)
	if err != nil {
		return nil, err
	}
	resp := &Response{
		Status:           stdResp.Status,
		StatusCode:       stdResp.StatusCode,
		Proto:            stdResp.Proto,
		ProtoMajor:       stdResp.ProtoMajor,
		ProtoMinor:       stdResp.ProtoMinor,
		Header:           Header(stdResp.Header),
		Body:             stdResp.Body,
		ContentLength:    stdResp.ContentLength,
		TransferEncoding: stdResp.TransferEncoding,
		Close:            stdResp.Close,
		Uncompressed:     stdResp.Uncompressed,
		Trailer:          Header(stdResp.Trailer),
		Request:          req,
		TLS:              stdResp.TLS,
	}
	return resp, nil
}
