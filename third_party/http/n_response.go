package http

import (
	"encoding/json"
	"io"
)

func (r *Response) readAndClose() []byte {
	body, _ := io.ReadAll(r.Body)
	_ = r.Body.Close()
	return body
}

func (r *Response) BytesBody() []byte {
	if r.cachedBody != nil {
		return r.cachedBody
	}

	body := r.readAndClose()
	r.cachedBody = body
	return body
}

func (r *Response) StringBody() string {
	return string(r.BytesBody())
}

func (r *Response) JsonBody(out interface{}) error {
	return json.Unmarshal(r.BytesBody(), out)
}
