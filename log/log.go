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

// Entry converts metadata fields into log.ContextLogger and injects it in
// context. By default all fields are passed to log context and multiple values
// are joined with ",".  If field names are specified, only those keys are
// extracted from metadata.
func Entry(ctx context.Context, fields ...string) context.Context {
	res := make(log.M)
	for k, v := range only(metadata.From(ctx), fields) {
		res[k] = strings.Join(v, ",")
	}
	return log.With(ctx, log.Get(ctx).WithContext(res))
}
