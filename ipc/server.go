package ipc

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	pb "github.com/zaker/oauth2local/ipc/localauth"
	"github.com/zaker/oauth2local/oauth2"
	"github.com/zaker/oauth2local/oauth2/redirect"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedLocalAuthServer
	oauthHandler oauth2.Handler
}

func NewServer(oauthHandler oauth2.Handler) (s *Server) {
	return &Server{oauthHandler: oauthHandler}
}

func (s *Server) GetAccessToken(ctx context.Context, _ *pb.Empty) (*pb.ATResponse, error) {
	slog.Debug("Get access token")
	a, err := s.oauthHandler.GetAccessToken()
	if err != nil {
		slog.Error("error:", "inner", err)
		return nil, err
	}
	return &pb.ATResponse{AccessToken: a}, nil
}

func (s *Server) Callback(ctx context.Context, cb *pb.CBRequest) (*pb.Empty, error) {
	slog.Debug("Callback from provider: ", "callback url", cb.Url)
	rURL, err := url.Parse(cb.Url)
	if err != nil {
		slog.Error("Callback error:", "callback error", err)
		return nil, err
	}
	err = s.oauthHandler.UpdateFromRedirect(redirect.DecodeRedirect(rURL))
	if err != nil {
		slog.Error("callback error:", "callback error", err)
		return nil, err
	}
	slog.Info("Auth ok")
	return &pb.Empty{}, nil
}

func (s *Server) Ping(ctx context.Context, _ *pb.Empty) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "pong"}, nil
}

func (s *Server) Serve() error {
	lis, err := listener()
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	defer lis.Close()

	gs := grpc.NewServer()
	pb.RegisterLocalAuthServer(gs, s)

	return gs.Serve(lis)
}
