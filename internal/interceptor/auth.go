package interceptor

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	secret      string
	unprotected map[string]bool
}

func NewAuthInterceptor(secret string, unprotected map[string]bool) *AuthInterceptor {
	return &AuthInterceptor{secret: secret, unprotected: unprotected}
}

func (a AuthInterceptor) Intercept(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if p := a.unprotected[info.FullMethod]; p == true {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	token := getTokenFromMetadata(md)
	if token == "" {
		return nil, fmt.Errorf("missing token")
	}

	user, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %s", err.Error())
	}

	ctx = context.WithValue(ctx, "user", user)
	return handler(ctx, req)
}

func getTokenFromMetadata(md metadata.MD) string {
	if authHeaders, ok := md["authorization"]; ok {
		if len(authHeaders) > 0 {
			return authHeaders[0]
		}
	}
	return ""
}
