package metadata

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-mixins/log"
	"github.com/go-mixins/metadata"
)

func only(md http.Header, fields []string) http.Header {
	if fields == nil {
		return md
	}
	res := make(http.Header, len(fields))
	for _, k := range fields {
		res[k] = md[k]
	}
	return res
}

// Fields converts metadata fields into map compatible with `go-mixins/log`. By
// default all fields are passed through and multiple values are joined with ",".
// If field names are specified, only those keys are extracted.
func Fields(ctx context.Context, fields ...string) log.M {
	res := make(log.M)
	for k, v := range only(metadata.From(ctx), fields) {
		res[k] = strings.Join(v, ",")
	}
	return res
}
