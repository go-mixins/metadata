// Package http provides HTTP client and server wrappers to send and receive
// metadata on the wire
package http

import (
	"net/http"

	"github.com/go-mixins/metadata"
)

// Handler wraps provided http.Handler and injects header fields into requests
// context metadata. All header names that start with `HeaderKeyPrefix` are
// have the prefix stripped and passed to request context automatically.
// The zero value of Handler is usable and wraps default HTTP server.
type Handler struct {
	// Handler is the handler used to handle the incoming request.
	Handler http.Handler
}

func (h *Handler) handler() http.Handler {
	if h.Handler != nil {
		return h.Handler
	}
	return http.DefaultServeMux
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := FromHeader(r.Context(), r.Header)
	h.handler().ServeHTTP(w, r.WithContext(ctx))
}

// Transport allows to pass metadata in outgoing HTTP requests. For
// compatibility with default request wrapper, all fields are converted to
// headers with `HeaderKeyPrefix` prepended to their names. The zero value is
// usable by default as a http.RoundTripper.
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
	req := r.WithContext(ctx) // shallow copy the request
	req.Header = metadata.Clone(req.Header)
	ToHeader(ctx, r.Header)
	return t.base().RoundTrip(req)
}
