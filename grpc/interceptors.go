package grpc

import (
	"context"
	"net/http"
	"net/textproto"
	"strings"

	"google.golang.org/grpc"
	grpcMetadata "google.golang.org/grpc/metadata"

	"github.com/go-mixins/metadata"
)

// MetadataKeyPrefix is prepended to metadata keys on the wire
var MetadataKeyPrefix = "x-grpc-meta-"

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md1 := metadata.From(ctx)
		md2, _ := grpcMetadata.FromOutgoingContext(ctx)
		dest := make(grpcMetadata.MD)
		for k, v := range md2 {
			vv := make([]string, len(v))
			copy(vv, v)
			dest[k] = vv
		}
		for k, v := range md1 {
			vv := make([]string, len(v))
			copy(vv, v)
			dest[MetadataKeyPrefix+strings.ToLower(k)] = vv
		}
		ctx = grpcMetadata.NewOutgoingContext(ctx, dest)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, rErr error) {
		md, _ := grpcMetadata.FromIncomingContext(ctx)
		dest := make(http.Header)
		for k, v := range md {
			if strings.HasPrefix(k, MetadataKeyPrefix) {
				dest[textproto.CanonicalMIMEHeaderKey(strings.TrimPrefix(k, MetadataKeyPrefix))] = v
			}
		}
		ctx = metadata.With(ctx, dest)
		return handler(ctx, req)
	}
}
