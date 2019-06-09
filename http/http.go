// Package http provides HTTP client and server wrappers to send and receive
// metadata on the wire
package http

import (
	"net/http"
	"strings"

	"github.com/go-mixins/metadata"
)

// Handler wrapws provided http.Handler and injects header fields into requests
// context metadata. All header names that start with `X-Meta-` are injected
// automatically. Specify `ExtraFields` to extract some extra fields
// from request header. The zero value of Handler is usable and wraps default
// HTTP server.
type Handler struct {
	// Handler is the handler used to handle the incoming request.
	Handler http.Handler
	// ExtraFields allows to pass arbitrary fields from incoming request
	// header to metadata. E.g.:
	//
	// &Handler{
	// 	ExtraFields: {"X-API-Key": "key"},
	// }
	ExtraFields map[string]string
}

func (h *Handler) handler() http.Handler {
	if h.Handler != nil {
		return h.Handler
	}
	return http.DefaultServeMux
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	md := make(http.Header)
	for k, vv := range r.Header {
		if newKey, ok := h.ExtraFields[k]; ok {
			vv2 := make([]string, len(vv))
			copy(vv2, vv)
			md[newKey] = vv
			continue
		}
		if newKey := strings.ToLower(k); strings.HasPrefix(newKey, "x-meta-") {
			vv2 := make([]string, len(vv))
			copy(vv2, vv)
			md[strings.TrimPrefix(newKey, "x-meta-")] = vv
		}
	}
	h.handler().ServeHTTP(w, r.WithContext(metadata.With(r.Context(), md)))
}

// Transport allows to pass metadata in outgoing HTTP requests. For
// compatibility with default request wrapper, all fields are converted to
// headers with `X-Meta-` name prefix. The zero value is usable by default as a
// http.RoundTripper.
type Transport struct {
	// Base may be set to wrap another http.RoundTripper
	Base http.RoundTripper
}

func (t *Transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

// CancelRequest cancels an in-flight request by closing its connection.
func (t *Transport) CancelRequest(req *http.Request) {
	type canceler interface {
		CancelRequest(*http.Request)
	}
	if cr, ok := t.base().(canceler); ok {
		cr.CancelRequest(req)
	}
}

// RoundTrip converts request context into headers
// and makes the request.
func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	ctx := r.Context()
	r = r.WithContext(ctx) // shallow copy the request
	for k, vv := range metadata.From(ctx) {
		r.Header[http.CanonicalHeaderKey("X-Meta-"+k)] = vv
	}
	return t.base().RoundTrip(r)
}
