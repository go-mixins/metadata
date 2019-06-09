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
// with `HeaderKeyPrefix` are lowercased and prefix is stripped. If field hame is
// present in `extraFields` map it is renamed to value for that key.
func FromHeader(src http.Header, extraFields map[string]string) http.Header {
	md := make(http.Header)
	for k, vv := range src {
		if newKey, ok := extraFields[k]; ok {
			vv2 := make([]string, len(vv))
			copy(vv2, vv)
			md[newKey] = vv
			continue
		} else if newKey := strings.ToLower(k); strings.HasPrefix(newKey, HeaderKeyPrefix) {
			vv2 := make([]string, len(vv))
			copy(vv2, vv)
			md[strings.TrimPrefix(newKey, HeaderKeyPrefix)] = vv
		}
	}
	return md
}

// ToHeader copies values from context metadata to specified http.Header,
// prepending `HeaderKeyPrefix` to field names.
func ToHeader(ctx context.Context, dest http.Header) {
	for k, vv := range metadata.From(ctx) {
		dest[http.CanonicalHeaderKey(HeaderKeyPrefix+k)] = vv
	}
}
