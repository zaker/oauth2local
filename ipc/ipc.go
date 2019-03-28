package ipc
import (
	"context"
	"fmt"
	"time"

	pb "github.com/equinor/oauth2local/ipc/localauth"
	"github.com/equinor/oauth2local/oauth2"
	"github.com/equinor/oauth2local/storage"
	"google.golang.org/grpc"
)

type Server struct {
	oauthCli oauth2.Client
	store    storage.Storage
}

type Client struct {
	grpcConn *grpc.ClientConn
	inner    pb.LocalAuthClient
}

func NewServer(cli oauth2.Client) (s *Server) {
	s = new(Server)
	s.oauthCli = cli
	return
}

func NewClient() (c *Client, err error) {
	c = new(Client)
	c.grpcConn, err = grpc.Dial("pipe", grpc.WithInsecure(), grpc.WithContextDialer(localPipeDial))
	if err != nil {
		return nil, err
	}
	c.inner = pb.NewLocalAuthClient(c.grpcConn)
	return
}

func (s *Server) GetAccessToken(ctx context.Context, _ *pb.Empty) (*pb.ATResponse, error) {
	r := new(pb.ATResponse)
	c, err := s.store.GetCode()
	a, err := s.oauthCli.GetToken(c)
	r.AccessToken = a
	return r, err
}

func (s *Server) UpdateCode(ctx context.Context, cr *pb.UCRequest) (*pb.Empty, error) {
	r := new(pb.Empty)
	fmt.Println("Received:", cr.Code)
	return r, nil

}

func (s *Server) Ping(ctx context.Context, _ *pb.Empty) (*pb.PingResponse, error) {
	r := new(pb.PingResponse)
	r.Message = "pong"
	return r, nil
}

func (c *Client) Close() {
	c.grpcConn.Close()
}

func HasSovereign() bool {
	c, err := NewClient()
	if err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.inner.Ping(ctx, new(pb.Empty))
	if err != nil {
		return false
	}
	return r.Message == "pong"
}

func (c *Client) SendCallback(callbackURL string) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	_, err := c.inner.UpdateCode(ctx, &pb.UCRequest{Code: callbackURL})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetAccessToken() (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	r, err := c.inner.GetAccessToken(ctx, &pb.Empty{})
	if err != nil {
		return "", err
	}
	return r.AccessToken, nil
}

func (s *Server) Serve() error {
	lis, err := listener()
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	gs := grpc.NewServer()
	pb.RegisterLocalAuthServer(gs, s)

	return gs.Serve(lis)
}
