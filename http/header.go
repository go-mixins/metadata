package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-mixins/metadata"
)

// HeaderKeyPrefix is prepended to meta keys when passing on the wire
var HeaderKeyPrefix = "x-meta-"

// FromHeader initializes new metadata from request header. All fields starting
// with `HeaderKeyPrefix` are lowercased and prefix is stripped. Optional list
// of extra header field names can be provided.
func FromHeader(ctx context.Context, src http.Header, extraFields ...string) context.Context {
	md := make(http.Header)
	for k, vv := range src {
		if newKey := strings.ToLower(k); strings.HasPrefix(newKey, HeaderKeyPrefix) {
			vv2 := make([]string, len(vv))
			copy(vv2, vv)
			md[strings.TrimPrefix(newKey, HeaderKeyPrefix)] = vv
		}
	}
	for _, k := range extraFields {
		vv := src[k]
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		md[k] = vv
	}
	return metadata.With(ctx, md)
}

// ToHeader copies values from context metadata to specified http.Header,
// prepending `HeaderKeyPrefix` to field names.
func ToHeader(ctx context.Context, dest http.Header) {
	for k, vv := range metadata.From(ctx) {
		dest[http.CanonicalHeaderKey(HeaderKeyPrefix+k)] = vv
	}
}
