package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-mixins/metadata"
)

// HeaderKeyPrefix is prepended to meta keys when passing on the wire
var HeaderKeyPrefix = "x-meta-"

// FromHeader merges metadata entries from request header. All fields starting
// with `HeaderKeyPrefix` are lowercased and prefix is stripped.
func FromHeader(ctx context.Context, src http.Header) context.Context {
	md := metadata.Clone(metadata.From(ctx))
	for k, vv := range src {
		if newKey := strings.ToLower(k); strings.HasPrefix(newKey, HeaderKeyPrefix) {
			vv2 := make([]string, len(vv))
			copy(vv2, vv)
			md[strings.TrimPrefix(newKey, HeaderKeyPrefix)] = vv
		}
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
