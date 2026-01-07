package services

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	gatewayv1 "github.com/whu-ham/ham-gateway/gen"
	"github.com/whu-ham/ham-gateway/internal/ham"
)

const (
	HeaderToken = "X-Gateway-Token"
)

type GatewayServer struct {
	hamClient *ham.Client
	authToken string
}

func NewGatewayServer(client *ham.Client, token string) *GatewayServer {
	return &GatewayServer{
		hamClient: client,
		authToken: token,
	}
}

func (s *GatewayServer) checkAuth(req connect.AnyRequest) error {
	if s.authToken == "" {
		return nil
	}
	token := req.Header().Get(HeaderToken)
	if token != s.authToken {
		return connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid token"))
	}
	return nil
}

func (s *GatewayServer) SearchCourse(
	ctx context.Context,
	req *connect.Request[gatewayv1.SearchCourseRequest],
) (*connect.Response[gatewayv1.SearchCourseResponse], error) {
	if err := s.checkAuth(req); err != nil {
		return nil, err
	}
	// Proxy to HAM Client
	// This is a placeholder for the actual mapping
	return connect.NewResponse(&gatewayv1.SearchCourseResponse{
		Code:    "00000",
		Message: "Success",
	}), nil
}

func (s *GatewayServer) GetCourseStat(
	ctx context.Context,
	req *connect.Request[gatewayv1.GetCourseStatRequest],
) (*connect.Response[gatewayv1.GetCourseStatResponse], error) {
	if err := s.checkAuth(req); err != nil {
		return nil, err
	}
	// Proxy to HAM Client
	return connect.NewResponse(&gatewayv1.GetCourseStatResponse{
		Code:    "00000",
		Message: "Success",
	}), nil
}
