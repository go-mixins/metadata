// Package metadata provides common utility functions to work with
// metadata-enabled context. It is intentionnaly mimics http.Header
// functionality and underlying data structure is also a http.Header.
package metadata

import (
	"context"
	"net/http"
)

type mdKey struct{}

// Clone header for thread safety
func Clone(h http.Header) http.Header {
	h2 := make(http.Header, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}

// With injects metadata into Context
func With(ctx context.Context, md http.Header) context.Context {
	return context.WithValue(ctx, mdKey{}, md)
}

// From extracts metadata from the Context. It is responsibility of user to
// clone it before modification.
func From(ctx context.Context) http.Header {
	md, _ := ctx.Value(mdKey{}).(http.Header)
	return md
}

// Set sets the metadata entries to single value
func Set(ctx context.Context, key, value string) context.Context {
	res := Clone(From(ctx))
	res.Set(key, value)
	return With(ctx, res)
}

// Add adds the key, value pair to the metadata
func Add(ctx context.Context, key, value string) context.Context {
	res := Clone(From(ctx))
	res.Add(key, value)
	return With(ctx, res)
}

// Del deletes the values associated with key from metadata
func Del(ctx context.Context, key string) context.Context {
	res := Clone(From(ctx))
	res.Del(key)
	return With(ctx, res)
}

// Get gets the first value associated with the given key
func Get(ctx context.Context, key string) string {
	md, _ := ctx.Value(mdKey{}).(http.Header)
	return md.Get(key)
}
