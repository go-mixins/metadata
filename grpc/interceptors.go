package grpc

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
	grpcMetadata "google.golang.org/grpc/metadata"

	"github.com/go-mixins/metadata"
)

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md := metadata.From(ctx)
		ctx = grpcMetadata.NewOutgoingContext(ctx, grpcMetadata.MD(md))
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, rErr error) {
		md, _ := grpcMetadata.FromIncomingContext(ctx)
		ctx = metadata.With(ctx, http.Header(md))
		return handler(ctx, req)
	}
}
