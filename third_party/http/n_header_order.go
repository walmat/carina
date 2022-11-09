package http

// HeaderOrderKey is a magic key for ResponseWriter.Header map keys
// that, if present, defines a header order that will be used to
// write the headers onto wire. The order of the list defined how the headers
// will be sorted. A defined key goes before an undefined key.
//
// This is the only way to specify some order, because maps don't
// have a a stable iteration order. If no order is given, headers will
// be sorted lexicographically.
//
// According to RFC-2616 it is good practice to send general-header fields
// first, followed by request-header or response-header fields and ending
// with entity-header fields.
const HeaderOrderKey = "Header-Order:"

func reqHasHeaderOrder(req *Request) bool {
	return req.Header.has(HeaderOrderKey)
}
