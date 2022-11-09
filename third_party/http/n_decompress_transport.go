package http

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"github.com/anknown/ahocorasick"
	"github.com/dsnet/compress/brotli"
	"io"
	"io/ioutil"
	"strings"
)

type DecompressTransport struct {
	inner RoundTripper
}

func NewDecompressTransport(inner RoundTripper) RoundTripper {
	return &DecompressTransport{inner: inner}
}

const (
	encodingGzip    = "gzip"
	encodingDeflate = "deflate"
	encodingBrotli  = "br"
)

const (
	zlibMethodDeflate = 0x78
	zlibLevelLow      = 0x01
	zlibLevelMedium   = 0x5E
	zlibLevelDefault  = 0x9C
	zlibLevelBest     = 0xDA
)

func (t *DecompressTransport) RoundTrip(req *Request) (*Response, error) {
	resp, err := t.inner.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	matcher := goahocorasick.Machine{}
	if err = matcher.Build([][]rune{[]rune(encodingGzip), []rune(encodingDeflate), []rune(encodingBrotli)}); err != nil {
		return resp, err
	}

	contentEncoding := strings.Join(resp.Header["Content-Encoding"], " ")
	encodingTerms := matcher.MultiPatternSearch([]rune(contentEncoding), false)
	lastEncodingIdx := len(encodingTerms) - 1

	isCompressed := false
	for idx := range encodingTerms {
		// we have to loop over it backwards because they are in reverse order in which they are applied (yes my comment earlier was wrong)
		encodingType := string(encodingTerms[lastEncodingIdx-idx].Word)
		switch encodingType {
		case encodingGzip:
			resp.Body = &nGzipReader{body: resp.Body}
			isCompressed = true
		case encodingDeflate:
			// TODO: find a way to not read entire body, we just need the first 2 bytes to check if it's zlib (proper) or the fucked deflate (thanks microsoft)
			respBody, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if err != nil {
				return nil, err
			}
			if respBody[0] == zlibMethodDeflate &&
				(respBody[1] == zlibLevelLow || respBody[1] == zlibLevelMedium || respBody[1] == zlibLevelDefault || respBody[1] == zlibLevelBest) {
				resp.Body = &nZlibReader{body: ioutil.NopCloser(bytes.NewReader(respBody))}
			} else {
				resp.Body = &nDeflateReader{body: ioutil.NopCloser(bytes.NewReader(respBody))}
			}
			isCompressed = true
		case encodingBrotli:
			resp.Body = &nBrotliReader{body: resp.Body}
			isCompressed = true
		}
	}

	if isCompressed {
		resp.Header.Del("Content-Encoding")
		resp.Header.Del("Content-Length")
		resp.ContentLength = -1
		resp.Uncompressed = true
	}

	return resp, nil
}

// nGzipReader wraps a response body so it can lazily
// call gzip.NewReader on the first call to Read
type nGzipReader struct {
	body io.ReadCloser // underlying Response.Body
	zr   *gzip.Reader  // lazily-initialized gzip reader
	zerr error         // sticky error
}

func (gz *nGzipReader) Read(p []byte) (n int, err error) {
	if gz.zerr != nil {
		return 0, gz.zerr
	}
	if gz.zr == nil {
		gz.zr, err = gzip.NewReader(gz.body)
		if err != nil {
			gz.zerr = err
			return 0, err
		}
	}
	return gz.zr.Read(p)
}

func (gz *nGzipReader) Close() error {
	return gz.body.Close()
}

// nBrotliReader wraps a response body so it can lazily
// call brotli.NewReader on the first call to Read
type nBrotliReader struct {
	body io.ReadCloser  // underlying Response.Body
	zr   *brotli.Reader // lazily-initialized brotli reader
	zerr error          // sticky error
}

func (gz *nBrotliReader) Read(p []byte) (n int, err error) {
	if gz.zerr != nil {
		return 0, gz.zerr
	}
	if gz.zr == nil {
		gz.zr, err = brotli.NewReader(gz.body, nil)
		if err != nil {
			gz.zerr = err
			return 0, err
		}
	}
	return gz.zr.Read(p)
}

func (gz *nBrotliReader) Close() error {
	return gz.body.Close()
}

// nDeflateReader wraps a response body so it can lazily
// call flate.NewReader on the first call to Read
type nDeflateReader struct {
	body io.ReadCloser // underlying Response.Body
	dr   io.ReadCloser // lazily-initialized deflate reader
	derr error         // sticky error
}

func (df *nDeflateReader) Read(p []byte) (n int, err error) {
	if df.dr == nil {
		df.dr = flate.NewReader(df.body)
	}
	return df.dr.Read(p)
}

func (df *nDeflateReader) Close() error {
	return df.body.Close()
}

// nZlibReader wraps a response body so it can lazily
// call zlib.NewReader on the first call to Read
type nZlibReader struct {
	body io.ReadCloser // underlying Response.Body
	zr   io.ReadCloser // lazily-initialized deflate reader
	zerr error         // sticky error
}

func (zr *nZlibReader) Read(p []byte) (n int, err error) {
	if zr.zerr != nil {
		return 0, zr.zerr
	}
	if zr.zr == nil {
		zr.zr, err = zlib.NewReader(zr.body)
		if err != nil {
			zr.zerr = err
			return 0, err
		}
	}
	return zr.zr.Read(p)
}

func (zr *nZlibReader) Close() error {
	return zr.body.Close()
}
