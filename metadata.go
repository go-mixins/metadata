// Package metadata provides common utility functions to work with
// metadata-enabled context. It is intentionnaly mimics http.Header
// functionality and underlying data structure is also a http.Header.
package metadata

import (
	"context"
	"net/http"
)

type m http.Header

type mdKey struct{}

func (h m) clone() http.Header {
	h2 := make(m, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return http.Header(h2)
}

// With injects metadata into Context
func With(ctx context.Context, md http.Header) context.Context {
	return context.WithValue(ctx, mdKey{}, md)
}

// From extracts copy of metadata from the Context
func From(ctx context.Context) http.Header {
	md, _ := ctx.Value(mdKey{}).(http.Header)
	return m(md).clone()
}

// Set sets the metadata entries to single value
func Set(ctx context.Context, key, value string) context.Context {
	res := From(ctx)
	res.Set(key, value)
	return context.WithValue(ctx, mdKey{}, res)
}

// Add adds the key, value pair to the metadata
func Add(ctx context.Context, key, value string) context.Context {
	res := From(ctx)
	res.Add(key, value)
	return context.WithValue(ctx, mdKey{}, res)
}

// Del deletes the values associated with key from metadata
func Del(ctx context.Context, key string) context.Context {
	res := From(ctx)
	res.Del(key)
	return With(ctx, res)
}

// Get gets the first value associated with the given key
func Get(ctx context.Context, key string) string {
	md, _ := ctx.Value(mdKey{}).(http.Header)
	return md.Get(key)
}
